package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type handler struct {
	hub        *notification.Hub
	activeRuns map[string]bool
	runsMu     sync.Mutex
}

func RegisterHandlers(
	r *gin.RouterGroup,
	hub *notification.Hub,
	authMiddleware gin.HandlerFunc,
) {
	h := &handler{
		hub:        hub,
		activeRuns: make(map[string]bool),
	}

	protected := r.Group("/notifications").Use(authMiddleware)
	protected.GET("/stream", h.stream)
}

func (h *handler) stream(c *gin.Context) {
	ctx := c.Request.Context()
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user id not found in context"})
		return
	}

	client := &notification.Client{
		UserID: userID,
		Send:   make(chan notification.Message, 16),
	}
	h.hub.Register(client)
	defer h.hub.Unregister(client)

	h.runsMu.Lock()
	if !h.activeRuns[userID] {
		h.activeRuns[userID] = true
		go func() {
			_ = h.hub.Run(context.Background(), userID)
			h.runsMu.Lock()
			delete(h.activeRuns, userID)
			h.runsMu.Unlock()
		}()
	}
	h.runsMu.Unlock()
	c.Header("Content-Type", "text/event-stream;charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	_, _ = fmt.Fprint(c.Writer, ": connected\n\n")
	c.Writer.Flush()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-client.Send:
			if !ok {
				return false
			}
			payload, err := json.Marshal(msg)
			if err != nil {
				logger.ErrorKV(ctx, "failed to marshal notification message", "error", err.Error(), "event_type", msg.EventType)
				return true
			}

			c.SSEvent(string(msg.EventType), string(payload))
			return true
		case <-ctx.Done():
			return false
		}
	})
}
