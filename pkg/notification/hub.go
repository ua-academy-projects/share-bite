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
	cancels      map[string]context.CancelFunc
}

func NewHub(s Subscriber) *Hub {
	return &Hub{
		clients:      make(map[string]map[*Client]bool),
		subscription: s,
		cancels:      make(map[string]context.CancelFunc),
	}
}

func (h *Hub) Register(c *Client) context.Context {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[c.UserID] == nil {
		h.clients[c.UserID] = make(map[*Client]bool)
	}
	h.clients[c.UserID][c] = true

	if _, exists := h.cancels[c.UserID]; !exists {
		ctx, cancel := context.WithCancel(context.Background())
		h.cancels[c.UserID] = cancel
		return ctx
	}

	return nil
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
				if cancel, hasCancel := h.cancels[c.UserID]; hasCancel {
					cancel()
					delete(h.cancels, c.UserID)
				}
			}
		}
	}
}

func (h *Hub) Run(ctx context.Context, channelName string) error {
	defer func() {
		h.mu.Lock()
		if ctx.Err() == nil {
			if cancel, hasCancel := h.cancels[channelName]; hasCancel {
				cancel()
				delete(h.cancels, channelName)
			}
		}
		h.mu.Unlock()
	}()

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
			userClients, exists := h.clients[msg.RecipientID]
			if exists {
				for client := range userClients {
					select {
					case client.Send <- msg:
					default:
						logger.DebugKV(ctx, "dropped notification", "user_id", msg.RecipientID, "client", client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}
