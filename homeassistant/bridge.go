package homeassistant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/timmo001/go-automate/config"
)

const defaultBridgeReconnectDelay = 5 * time.Second

type Bridge struct {
	cfg            *config.ConfigHomeAssistant
	socketPath     string
	reconnectDelay time.Duration

	stateMu sync.RWMutex
	states  map[string]HomeAssistantState

	subscriberMu sync.Mutex
	nextSubID    uint64
	subscribers  map[string]map[uint64]chan HomeAssistantState
}

type BridgeRequest struct {
	Action   string `json:"action"`
	EntityID string `json:"entity_id,omitempty"`
}

type BridgeResponse struct {
	Type     string              `json:"type"`
	EntityID string              `json:"entity_id,omitempty"`
	State    *HomeAssistantState `json:"state,omitempty"`
	Error    string              `json:"error,omitempty"`
	Meta     map[string]string   `json:"meta,omitempty"`
}

func DefaultBridgeSocketPath() (string, error) {
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if runtimeDir == "" {
		runtimeDir = os.TempDir()
	}

	if runtimeDir == "" {
		return "", fmt.Errorf("could not determine runtime directory")
	}

	return filepath.Join(runtimeDir, "go-automate", "home-assistant.sock"), nil
}

func NewBridge(cfg *config.ConfigHomeAssistant, socketPath string) (*Bridge, error) {
	if cfg == nil {
		return nil, fmt.Errorf("home assistant config is not set")
	}

	if socketPath == "" {
		var err error
		socketPath, err = DefaultBridgeSocketPath()
		if err != nil {
			return nil, err
		}
	}

	return &Bridge{
		cfg:            cfg,
		socketPath:     socketPath,
		reconnectDelay: defaultBridgeReconnectDelay,
		states:         map[string]HomeAssistantState{},
		subscribers:    map[string]map[uint64]chan HomeAssistantState{},
	}, nil
}

func (bridge *Bridge) SocketPath() string {
	return bridge.socketPath
}

func (bridge *Bridge) Serve(ctx context.Context) error {
	if err := os.MkdirAll(filepath.Dir(bridge.socketPath), 0700); err != nil {
		return fmt.Errorf("create bridge socket directory: %w", err)
	}

	if err := os.Remove(bridge.socketPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove stale bridge socket: %w", err)
	}

	listener, err := net.Listen("unix", bridge.socketPath)
	if err != nil {
		return fmt.Errorf("listen on bridge socket: %w", err)
	}
	defer listener.Close()
	defer os.Remove(bridge.socketPath)

	if err := os.Chmod(bridge.socketPath, 0600); err != nil {
		return fmt.Errorf("set bridge socket permissions: %w", err)
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	go bridge.runUpstream(ctx)

	log.Infof("Home Assistant bridge listening on %s", bridge.socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			log.Errorf("Error accepting bridge connection: %v", err)
			continue
		}

		go bridge.handleClient(ctx, conn)
	}
}

func (bridge *Bridge) runUpstream(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		conn, err := Dial(bridge.cfg)
		if err != nil {
			log.Errorf("Error connecting bridge to Home Assistant: %v", err)
			if !bridge.waitForReconnect(ctx) {
				return
			}
			continue
		}

		if err := bridge.syncStates(conn); err != nil {
			log.Errorf("Error syncing Home Assistant state: %v", err)
			conn.Close()
			if !bridge.waitForReconnect(ctx) {
				return
			}
			continue
		}

		response, err := conn.SubscribeEvents("state_changed")
		if err != nil {
			log.Errorf("Error subscribing to Home Assistant events: %v", err)
			conn.Close()
			if !bridge.waitForReconnect(ctx) {
				return
			}
			continue
		}
		if !response.Success {
			log.Errorf("Home Assistant subscribe failed: %s", response.Error.Message)
			conn.Close()
			if !bridge.waitForReconnect(ctx) {
				return
			}
			continue
		}

		log.Info("Home Assistant bridge subscribed to state_changed")

		for {
			if ctx.Err() != nil {
				conn.Close()
				return
			}

			event, err := conn.ReadEvent()
			if err != nil {
				log.Errorf("Error reading Home Assistant event: %v", err)
				conn.Close()
				break
			}

			if event.Type != "event" || event.Event.EventType != "state_changed" || event.Event.Data.NewState == nil {
				continue
			}

			bridge.storeState(*event.Event.Data.NewState)
		}

		if !bridge.waitForReconnect(ctx) {
			return
		}
	}
}

func (bridge *Bridge) waitForReconnect(ctx context.Context) bool {
	timer := time.NewTimer(bridge.reconnectDelay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (bridge *Bridge) syncStates(conn *HomeAssistantConn) error {
	states, err := conn.GetStates()
	if err != nil {
		return err
	}

	bridge.stateMu.Lock()
	bridge.states = make(map[string]HomeAssistantState, len(states))
	for _, state := range states {
		bridge.states[state.EntityID] = state
	}
	bridge.stateMu.Unlock()

	bridge.subscriberMu.Lock()
	entityIDs := make([]string, 0, len(bridge.subscribers))
	for entityID := range bridge.subscribers {
		entityIDs = append(entityIDs, entityID)
	}
	bridge.subscriberMu.Unlock()

	for _, entityID := range entityIDs {
		if state, ok := bridge.getState(entityID); ok {
			bridge.broadcastState(state)
		}
	}

	return nil
}

func (bridge *Bridge) getState(entityID string) (HomeAssistantState, bool) {
	bridge.stateMu.RLock()
	defer bridge.stateMu.RUnlock()

	state, ok := bridge.states[entityID]
	return state, ok
}

func (bridge *Bridge) storeState(state HomeAssistantState) {
	bridge.stateMu.Lock()
	bridge.states[state.EntityID] = state
	bridge.stateMu.Unlock()

	bridge.broadcastState(state)
}

func (bridge *Bridge) broadcastState(state HomeAssistantState) {
	bridge.subscriberMu.Lock()
	subscribers := bridge.subscribers[state.EntityID]
	channels := make([]chan HomeAssistantState, 0, len(subscribers))
	for _, ch := range subscribers {
		channels = append(channels, ch)
	}
	bridge.subscriberMu.Unlock()

	for _, ch := range channels {
		select {
		case ch <- state:
		default:
		}
	}
}

func (bridge *Bridge) addSubscriber(entityID string) (uint64, chan HomeAssistantState) {
	bridge.subscriberMu.Lock()
	defer bridge.subscriberMu.Unlock()

	bridge.nextSubID++
	subID := bridge.nextSubID
	ch := make(chan HomeAssistantState, 16)

	if bridge.subscribers[entityID] == nil {
		bridge.subscribers[entityID] = map[uint64]chan HomeAssistantState{}
	}
	bridge.subscribers[entityID][subID] = ch

	return subID, ch
}

func (bridge *Bridge) removeSubscriber(entityID string, subID uint64) {
	bridge.subscriberMu.Lock()
	defer bridge.subscriberMu.Unlock()

	entitySubscribers := bridge.subscribers[entityID]
	if entitySubscribers == nil {
		return
	}

	delete(entitySubscribers, subID)
	if len(entitySubscribers) == 0 {
		delete(bridge.subscribers, entityID)
	}
}

func (bridge *Bridge) handleClient(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	var request BridgeRequest
	if err := decoder.Decode(&request); err != nil {
		bridge.writeBridgeError(encoder, fmt.Errorf("decode bridge request: %w", err))
		return
	}

	switch request.Action {
	case "get_entity":
		if request.EntityID == "" {
			bridge.writeBridgeError(encoder, fmt.Errorf("entity_id is required"))
			return
		}

		state, ok := bridge.getState(request.EntityID)
		if !ok {
			if err := encoder.Encode(BridgeResponse{
				Type:     "snapshot",
				EntityID: request.EntityID,
				State:    nil,
			}); err != nil {
				log.Errorf("Error writing bridge snapshot response: %v", err)
			}
			return
		}

		stateCopy := state
		if err := encoder.Encode(BridgeResponse{
			Type:     "snapshot",
			EntityID: request.EntityID,
			State:    &stateCopy,
		}); err != nil {
			log.Errorf("Error writing bridge snapshot response: %v", err)
		}
	case "watch_entity":
		if request.EntityID == "" {
			bridge.writeBridgeError(encoder, fmt.Errorf("entity_id is required"))
			return
		}

		subID, stateCh := bridge.addSubscriber(request.EntityID)
		defer bridge.removeSubscriber(request.EntityID, subID)

		if state, ok := bridge.getState(request.EntityID); ok {
			stateCopy := state
			if err := encoder.Encode(BridgeResponse{
				Type:     "snapshot",
				EntityID: request.EntityID,
				State:    &stateCopy,
			}); err != nil {
				log.Errorf("Error writing bridge snapshot response: %v", err)
				return
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case state := <-stateCh:
				stateCopy := state
				if err := encoder.Encode(BridgeResponse{
					Type:     "state_changed",
					EntityID: request.EntityID,
					State:    &stateCopy,
				}); err != nil {
					log.Errorf("Error writing bridge event response: %v", err)
					return
				}
			}
		}
	default:
		bridge.writeBridgeError(encoder, fmt.Errorf("unsupported bridge action %q", request.Action))
	}
}

func (bridge *Bridge) writeBridgeError(encoder *json.Encoder, err error) {
	if encodeErr := encoder.Encode(BridgeResponse{
		Type:  "error",
		Error: err.Error(),
	}); encodeErr != nil {
		log.Errorf("Error writing bridge error response: %v", encodeErr)
	}
}

func BridgeWatchEntity(
	ctx context.Context,
	socketPath string,
	entityID string,
	onState func(*HomeAssistantState) error,
) error {
	if socketPath == "" {
		var err error
		socketPath, err = DefaultBridgeSocketPath()
		if err != nil {
			return err
		}
	}

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return fmt.Errorf("connect to bridge socket: %w", err)
	}
	defer conn.Close()

	if err := json.NewEncoder(conn).Encode(BridgeRequest{
		Action:   "watch_entity",
		EntityID: entityID,
	}); err != nil {
		return fmt.Errorf("send bridge watch request: %w", err)
	}

	type result struct {
		err error
	}

	resultCh := make(chan result, 1)
	go func() {
		decoder := json.NewDecoder(conn)
		for {
			var response BridgeResponse
			if err := decoder.Decode(&response); err != nil {
				resultCh <- result{err: fmt.Errorf("read bridge response: %w", err)}
				return
			}

			if response.Type == "error" {
				resultCh <- result{err: fmt.Errorf("bridge error: %s", response.Error)}
				return
			}

			if response.State == nil {
				continue
			}

			if err := onState(response.State); err != nil {
				resultCh <- result{err: err}
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case result := <-resultCh:
		return result.err
	}
}
