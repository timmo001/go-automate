package homeassistant

import (
	"encoding/json"
	"net/url"
	"strings"

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

type HomeAssistantCallServiceRequest struct {
	ID             int         `json:"id"`
	Type           string      `json:"type"`
	Domain         string      `json:"domain"`
	Service        string      `json:"service"`
	ServiceData    interface{} `json:"service_data,omitempty"`
	Target         interface{} `json:"target,omitempty"`
	ReturnResponse bool        `json:"return_response,omitempty"`
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
	Config *config.ConfigHomeAssistant
)

func Connect() *HomeAssistantConn {
	// 1. Setup WebSocket URL
	host := strings.Split(Config.URL, "/")[2]

	wsUrl := url.URL{Scheme: "ws", Host: host, Path: "/api/websocket"}
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
