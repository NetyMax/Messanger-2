package ws

import (
	"sync"
)

type Hub struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var hub = Hub{
	clients: make(map[string]*Client),
}

func RegisterClient(userID string, client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.clients[userID] = client
}

func UnRegisterCkient(userID string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	delete(hub.clients, userID)
}

func SendToUser(to string, message []byte) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	if client, ok := hub.clients[to]; ok {
		client.send <- message
	}
}

func HubMu() *sync.RWMutex {
	return &hub.mu
}

func HubClients() map[string]*Client {
	return hub.clients
}
