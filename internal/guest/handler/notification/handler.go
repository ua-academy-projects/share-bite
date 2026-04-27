package notification

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type handler struct {
	hub        *notification.Hub
	customer   customerService
	activeRuns map[string]bool
	runsMu     sync.Mutex
}

type customerService interface {
	GetByUserID(ctx context.Context, userID string) (entity.Customer, error)
}

func RegisterHandlers(
	r *gin.RouterGroup,
	hub *notification.Hub,
	customerService customerService,
	authMiddleware gin.HandlerFunc,
) {
	h := &handler{
		hub:        hub,
		customer:   customerService,
		activeRuns: make(map[string]bool),
	}
	protected := r.Group("/").Use(authMiddleware)

	protected.GET("/stream", h.stream)
}

func (h *handler) stream(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if h.customer != nil {
		if _, err := h.customer.GetByUserID(ctx, userID); err != nil {
			c.Error(err)
			return
		}
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

			c.SSEvent(string(msg.Type), msg.Data)
			return true
		case <-ctx.Done():
			return false
		}
	})
}
