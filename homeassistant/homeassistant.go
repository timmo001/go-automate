package homeassistant

import (
	"encoding/json"
	"net/url"
	"sync/atomic"

	"github.com/timmo001/go-automate/config"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

type HomeAssistantConn struct {
	*websocket.Conn
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
	parsedURL, err := url.Parse(Config.URL)
	if err != nil {
		log.Fatalf("Error parsing Home Assistant URL: %v", err)
	}

	wsScheme := "ws"
	if parsedURL.Scheme == "https" || parsedURL.Scheme == "wss" {
		wsScheme = "wss"
	}

	wsUrl := url.URL{Scheme: wsScheme, Host: parsedURL.Host, Path: "/api/websocket"}
	log.Infof("Connecting to Home Assistant at: %s", wsUrl.String())

	// 2. Connect to WebSocket
	c, _, err := websocket.DefaultDialer.Dial(wsUrl.String(), nil)
	if err != nil {
		log.Fatalf("Error connecting to WebSocket: %v", err)
	}
	log.Info("Connected to WebSocket")
	conn := &HomeAssistantConn{Conn: c}

	// 3. Read welcome message
	var welcomeResponse HomeAssistantResponse[any]
	err = conn.ReadJSON(&welcomeResponse)
	if err != nil {
		log.Fatalf("Error reading welcome message: %v", err)
	}
	log.Infof("First response: %v", welcomeResponse)

	// 4. Create auth message
	authMessage := HomeAssistantAuthRequest{
		Type:        "auth",
		AccessToken: Config.Token,
	}

	// 5. Send auth message and read response
	authResponse := conn.SendRequest(authMessage, false)
	log.Infof("Auth response: %v", authResponse)

	return conn
}

func RandomID() int {
	return int(requestIDSeed.Add(1))
}

func (conn *HomeAssistantConn) SendRequest(request interface{}, debug bool) HomeAssistantResponse[any] {
	if debug {
		b, err := json.Marshal(request)
		if err != nil {
			log.Errorf("Error marshalling request: %v", err)
		} else {
			log.Infof("Request: %s", b)
		}
	}
	err := conn.WriteJSON(request)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	var response HomeAssistantResponse[any]
	err = conn.ReadJSON(&response)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	return response
}

func (conn *HomeAssistantConn) GetState(entityID string) *HomeAssistantState {
	request := HomeAssistantRequest{
		ID:   RandomID(),
		Type: "get_states",
	}

	b, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("Error marshalling get_states request: %v", err)
	}
	log.Infof("Request: %s", b)

	err = conn.WriteJSON(request)
	if err != nil {
		log.Fatalf("Error sending get_states request: %v", err)
	}

	var response HomeAssistantResponse[[]HomeAssistantState]
	err = conn.ReadJSON(&response)
	if err != nil {
		log.Fatalf("Error reading get_states response: %v", err)
	}

	for _, state := range response.Result {
		if state.EntityID == entityID {
			state := state
			return &state
		}
	}

	return nil
}

func (conn *HomeAssistantConn) SubscribeEvents(eventType string) HomeAssistantResponse[any] {
	return conn.SendRequest(HomeAssistantSubscribeEventsRequest{
		ID:        RandomID(),
		Type:      "subscribe_events",
		EventType: eventType,
	}, true)
}

func (conn *HomeAssistantConn) ReadEvent() HomeAssistantEventMessage {
	var event HomeAssistantEventMessage
	err := conn.ReadJSON(&event)
	if err != nil {
		log.Fatalf("Error reading event: %v", err)
	}

	return event
}
