package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Max/messenger-2/ws"
	"github.com/gorilla/websocket"
)

func TestHandleWebSocket_NoUserID(t *testing.T) {
	req := httptest.NewRequest("GET", "/ws", nil)
	rr := httptest.NewRecorder()

	ws.HandleWebSocket(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", rr.Code)
	}
}

func TestHandleWebSocket_ValidRequest(t *testing.T) {
	// Создаём HTTP-сервер
	s := httptest.NewServer(http.HandlerFunc(ws.HandleWebSocket))
	defer s.Close()

	// Преобразуем адрес в ws:// вместо http://
	u := "ws" + strings.TrimPrefix(s.URL, "http") + "/?user=testuser"

	// Устанавливаем WebSocket-соединение
	wsConn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}
	defer wsConn.Close()
}

func TestClientReadAndWrite(t *testing.T) {
	var receivedFrom string
	var receivedMsg []byte

	ws.HandleMessage = func(from string, msg []byte) {
		receivedFrom = from
		receivedMsg = msg
	}

	server := httptest.NewServer(http.HandleFunc(ws.HandleWebSocket))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/?user=testuser"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}

	time.Sleep(100 * time.Millisecond)
	if receivedFrom != "testuser" {
		t.Errorf("expected from = testuser, got %s", receivedFrom)

	}

	if string(receivedMsg) == "" {
		t.Errorf("expected non-empty message")
	}

	ws.SendToUser("testuser", []byte("echo test"))

	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("rean error: %v", err)
	}

	if string(msg) != "echo test" {
		t.Errorf("Expected 'echo test', got %s", msg)
	}
}

var sentTo string
var sentMsg []byte

func init() {
	ws.SendToUser = func(to string, message []byte) {
		sentTo = to
		sentMsg = message
	}
}

func TestHandleMessage_Valid(t *testing.T) {
	input := ws.Message{
		Type:    "message",
		To:      "user2",
		Content: "hello",
	}

	data, _ := json.Marshal(input)

	ws.HandleMessage("user1", data)

	if sentTo != "user2" {
		t.Errorf("expected message to user2", sentTo)
	}

	var parsed struct {
		From    string `json:"from"`
		Content string `json:"content"`
	}

	err := json.Unmarshal(sentMsg, &parsed)
	if err != nil {
		t.Fatalf("INVALID", err)
	}

	if parsed.From != "user1" || parsed.Content != "hello" {
		t.Errorf("expected from")
	}
}

func TestRegisterAndUnregisterClient(t *testing.T) {
	client := &ws.Client{send: make(chan []byte, 1)}
	ws.RegisterClient("user1", client)

	ws.HubMu().RLock()
	if ws.HubClients()["user1"] != client {
		t.Errorf("error")
	}

	ws.HubMu().RUnlock()
	ws.RegisterClient("user1", client)
	ws.HubMu().RLock()
	if _, ok := ws.HubClients()["user1"]; ok {
		t.Errorf("error")
	}

	ws.HubMu().RUnlock()
}

func TestSendToUser(t *testing.T) {
	message := []byte("hello")

	recvChan := make(chan []byte, 1)
	client := &ws.Client{send: recvChan}
	ws.RegisterClient("user2", client)

	ws.SendToUser("user2", message)

	select {
	case msg := <-recvChan:
		if string(msg) != "hello" {
			t.Errorf("error", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("error")
	}

	ws.UnRegisterClient("user2")
}
