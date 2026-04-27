package notification

import (
	"context"
	"sync"

	"github.com/ua-academy-projects/share-bite/pkg/logger"
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
	if userClients, ok := h.clients[c.UserID]; ok {
		if _, exists := userClients[c]; exists {
			delete(userClients, c)
			close(c.Send)
			if len(userClients) == 0 {
				delete(h.clients, c.UserID)
			}
		}
	}
}

func (h *Hub) Run(ctx context.Context, channelName string) error {
	msgChan, err := h.subscription.Subscribe(ctx, channelName)
	if err != nil {
		logger.ErrorKV(ctx, "failed to subscribe to notification channel", "channel", channelName, "error", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}

			h.mu.RLock()
			userClients, exists := h.clients[msg.UserID]
			if exists {
				for client := range userClients {
					select {
					case client.Send <- msg:
					default:
						logger.DebugKV(ctx, "dropped notification", "user_id", msg.UserID, "client", client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}
