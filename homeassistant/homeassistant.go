package homeassistant

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/timmo001/go-automate/config"
)

type HomeAssistantConn struct {
	*websocket.Conn
	writeMu sync.Mutex
}

type HomeAssistantRequest struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

type HomeAssistantAuthRequest struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type HomeAssistantSubscribeEventsRequest struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type,omitempty"`
}

type HomeAssistantCallServiceRequest struct {
	ID             int         `json:"id"`
	Type           string      `json:"type"`
	Domain         string      `json:"domain"`
	Service        string      `json:"service"`
	ServiceData    interface{} `json:"service_data,omitempty"`
	Target         interface{} `json:"target,omitempty"`
	ReturnResponse bool        `json:"return_response,omitempty"`
}

type HomeAssistantState struct {
	EntityID string `json:"entity_id"`
	State    string `json:"state"`
}

type HomeAssistantEventMessage struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Event struct {
		EventType string `json:"event_type"`
		Data      struct {
			EntityID string              `json:"entity_id"`
			OldState *HomeAssistantState `json:"old_state"`
			NewState *HomeAssistantState `json:"new_state"`
		} `json:"data"`
	} `json:"event"`
}

type HomeAssistantResponse[T any] struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result    T      `json:"result"`
	HaVersion string `json:"ha_version"`
}

var (
	Config        *config.ConfigHomeAssistant
	requestIDSeed atomic.Int64
)

func Connect() *HomeAssistantConn {
	conn, err := Dial(Config)
	if err != nil {
		log.Fatalf("Error connecting to Home Assistant: %v", err)
	}

	return conn
}

func Dial(cfg *config.ConfigHomeAssistant) (*HomeAssistantConn, error) {
	if cfg == nil {
		return nil, fmt.Errorf("home assistant config is not set")
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("home assistant URL is not set")
	}

	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse home assistant URL: %w", err)
	}

	wsScheme := "ws"
	if parsedURL.Scheme == "https" || parsedURL.Scheme == "wss" {
		wsScheme = "wss"
	}

	wsUrl := url.URL{Scheme: wsScheme, Host: parsedURL.Host, Path: "/api/websocket"}
	log.Infof("Connecting to Home Assistant at: %s", wsUrl.String())

	c, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("connect to Home Assistant websocket: %w", err)
	}
	log.Info("Connected to WebSocket")
	conn := &HomeAssistantConn{Conn: c}

	if err := conn.authenticate(cfg.Token); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func RandomID() int {
	return int(requestIDSeed.Add(1))
}

func (conn *HomeAssistantConn) authenticate(token string) error {
	var welcomeResponse HomeAssistantResponse[any]
	if err := conn.ReadJSON(&welcomeResponse); err != nil {
		return fmt.Errorf("read welcome message: %w", err)
	}
	log.Infof("First response: %v", welcomeResponse)

	authMessage := HomeAssistantAuthRequest{
		Type:        "auth",
		AccessToken: token,
	}

	var authResponse HomeAssistantResponse[any]
	if err := conn.sendRequestInto(authMessage, false, &authResponse); err != nil {
		return fmt.Errorf("authenticate websocket: %w", err)
	}
	log.Infof("Auth response: %v", authResponse)

	if authResponse.Type != "auth_ok" {
		return fmt.Errorf("authentication failed with response type %q", authResponse.Type)
	}

	return nil
}

func (conn *HomeAssistantConn) SendRequest(request any, debug bool) (HomeAssistantResponse[any], error) {
	var response HomeAssistantResponse[any]
	if err := conn.sendRequestInto(request, debug, &response); err != nil {
		return HomeAssistantResponse[any]{}, err
	}

	return response, nil
}

func (conn *HomeAssistantConn) sendRequestInto(request any, debug bool, response any) error {
	if debug {
		b, err := json.Marshal(request)
		if err != nil {
			log.Errorf("Error marshalling request: %v", err)
		} else {
			log.Infof("Request: %s", b)
		}
	}
	conn.writeMu.Lock()
	err := conn.WriteJSON(request)
	conn.writeMu.Unlock()
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	if err := conn.ReadJSON(response); err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	return nil
}

func (conn *HomeAssistantConn) GetStates() ([]HomeAssistantState, error) {
	request := HomeAssistantRequest{
		ID:   RandomID(),
		Type: "get_states",
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal get_states request: %w", err)
	}
	log.Infof("Request: %s", b)

	var response HomeAssistantResponse[[]HomeAssistantState]
	if err := conn.sendRequestInto(request, false, &response); err != nil {
		return nil, fmt.Errorf("get states: %w", err)
	}

	return response.Result, nil
}

func (conn *HomeAssistantConn) GetState(entityID string) (*HomeAssistantState, error) {
	states, err := conn.GetStates()
	if err != nil {
		return nil, err
	}

	for _, state := range states {
		if state.EntityID == entityID {
			state := state
			return &state, nil
		}
	}

	return nil, nil
}

func (conn *HomeAssistantConn) SubscribeEvents(eventType string) (HomeAssistantResponse[any], error) {
	return conn.SendRequest(HomeAssistantSubscribeEventsRequest{
		ID:        RandomID(),
		Type:      "subscribe_events",
		EventType: eventType,
	}, true)
}

func (conn *HomeAssistantConn) ReadEvent() (HomeAssistantEventMessage, error) {
	var event HomeAssistantEventMessage
	if err := conn.ReadJSON(&event); err != nil {
		return HomeAssistantEventMessage{}, fmt.Errorf("read event: %w", err)
	}

	return event, nil
}
