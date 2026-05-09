package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	guestentity "github.com/ua-academy-projects/share-bite/internal/guest/entity"
	"github.com/ua-academy-projects/share-bite/internal/notification/service"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	notificationpkg "github.com/ua-academy-projects/share-bite/pkg/notification"
)

type customerService interface {
	GetByUserID(ctx context.Context, userID string) (guestentity.Customer, error)
}

type handler struct {
	svc        *service.Service
	hub        *notificationpkg.Hub
	customers  customerService
	activeRuns map[string]bool
	runsMu     sync.Mutex
}

func RegisterHandlers(r *gin.RouterGroup, svc *service.Service, hub *notificationpkg.Hub, customers customerService, authMiddleware gin.HandlerFunc, streamAuthMiddleware gin.HandlerFunc) {
	h := &handler{
		svc:        svc,
		hub:        hub,
		customers:  customers,
		activeRuns: make(map[string]bool),
	}

	protected := r.Group("/").Use(authMiddleware)
	protected.GET("/", h.getHistory)
	protected.POST("/mark-read", h.markAsRead)
	stream := r.Group("/").Use(streamAuthMiddleware)
	stream.GET("/stream", h.stream)
}

type notificationResponse struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	EntityID  string         `json:"entity_id"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

func (h *handler) getHistory(c *gin.Context) {
	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if h.customers != nil {
		if _, err := h.customers.GetByUserID(c.Request.Context(), userID); err != nil {
			c.Error(err)
			return
		}
	}

	limit := 20
	if raw := c.Query("limit"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	offset := 0
	if raw := c.Query("offset"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	items, err := h.svc.GetHistory(c.Request.Context(), userID, limit, offset)
	if err != nil {
		logger.ErrorKV(c.Request.Context(), "get notification history", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notifications"})
		return
	}

	response := make([]notificationResponse, 0, len(items))
	for _, item := range items {
		response = append(response, notificationResponse{
			ID:        item.ID,
			Type:      item.Type,
			EntityID:  item.EntityID,
			Metadata:  item.Metadata,
			CreatedAt: item.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) markAsRead(c *gin.Context) {
	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req struct {
		NotificationIDs []string `json:"notification_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.svc.MarkAsRead(c.Request.Context(), userID, req.NotificationIDs); err != nil {
		logger.ErrorKV(c.Request.Context(), "mark notifications as read", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark notifications as read"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *handler) stream(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := httpctx.GetUserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if h.customers != nil {
		if _, err := h.customers.GetByUserID(ctx, userID); err != nil {
			c.Error(err)
			return
		}
	}

	client := &notificationpkg.Client{
		UserID: userID,
		Send:   make(chan notificationpkg.Message, 16),
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
			c.SSEvent(string(msg.EventType), notificationResponse{
				ID:        msg.EventID,
				Type:      string(msg.EventType),
				EntityID:  msg.EntityID,
				Metadata:  msg.Metadata,
				CreatedAt: msg.CreatedAt,
			})
			return true
		case <-ctx.Done():
			return false
		}
	})
}
