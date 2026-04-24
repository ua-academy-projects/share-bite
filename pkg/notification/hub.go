package notification

import (
	"context"
	"sync"
)

type Client struct {
	UserID string
	Send   chan Message
}
type Hub struct {
	mu           sync.RWMutex
	clients      map[string]map[*Client]bool
	subscription Subscriber
}

func NewHub(s Subscriber) *Hub {
	return &Hub{
		clients:      make(map[string]map[*Client]bool),
		subscription: s,
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[c.UserID] == nil {
		h.clients[c.UserID] = make(map[*Client]bool)
	}
	h.clients[c.UserID][c] = true
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[c.UserID]; ok {
		delete(h.clients[c.UserID], c)
		close(c.Send)
		if len(h.clients[c.UserID]) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

func (h *Hub) Run(ctx context.Context, channelName string) {
	msgChan, err := h.subscription.Subscribe(ctx, channelName)
	if err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgChan:
			if !ok {
				return
			}

			h.mu.RLock()
			userClients, exists := h.clients[msg.UserID]
			if exists {
				for client := range userClients {
					select {
					case client.Send <- msg:
					default:
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}
