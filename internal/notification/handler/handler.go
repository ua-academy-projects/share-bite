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
	"github.com/ua-academy-projects/share-bite/internal/notification/service"
	"github.com/ua-academy-projects/share-bite/internal/util/httpctx"
	"github.com/ua-academy-projects/share-bite/internal/util/request"
	"github.com/ua-academy-projects/share-bite/pkg/logger"
	notificationpkg "github.com/ua-academy-projects/share-bite/pkg/notification"
)

type handler struct {
	svc        *service.Service
	hub        *notificationpkg.Hub
	userSubs   map[string]int
	userCancel map[string]context.CancelFunc
	runsMu     sync.Mutex
}

func RegisterHandlers(r *gin.RouterGroup, svc *service.Service, hub *notificationpkg.Hub, authMiddleware gin.HandlerFunc, streamAuthMiddleware gin.HandlerFunc) {
	h := &handler{
		svc:        svc,
		hub:        hub,
		userSubs:   make(map[string]int),
		userCancel: make(map[string]context.CancelFunc),
	}

	auth := r.Group("/").Use(authMiddleware)
	{
		auth.GET("/history", h.getHistory)
		auth.POST("/mark-read", h.markAsRead)
	}

	stream := r.Group("/").Use(streamAuthMiddleware)
	{
		stream.GET("/stream", h.stream)
	}
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
		c.Error(err)
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
	if err := request.BindJSON(c, &req); err != nil {
		c.Error(err)
		return
	}

	if err := h.svc.MarkAsRead(c.Request.Context(), userID, req.NotificationIDs); err != nil {
		logger.ErrorKV(c.Request.Context(), "mark notifications as read", "error", err)
		c.Error(err)
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

	logger.InfoKV(ctx, "new sse connection", "user_id", userID)

	client := &notificationpkg.Client{
		UserID: userID,
		Send:   make(chan notificationpkg.Message, 16),
	}
	h.hub.Register(client)
	h.runsMu.Lock()
	h.userSubs[userID]++
	if h.userSubs[userID] == 1 {
		ctxRun, cancel := context.WithCancel(context.Background())
		h.userCancel[userID] = cancel
		go func() {
			_ = h.hub.Run(ctxRun, userID)
		}()
	}
	h.runsMu.Unlock()

	defer func() {
		h.hub.Unregister(client)
		h.runsMu.Lock()
		h.userSubs[userID]--
		if h.userSubs[userID] == 0 {
			if cancel, ok := h.userCancel[userID]; ok {
				cancel()
				delete(h.userCancel, userID)
			}
			delete(h.userSubs, userID)
		}
		h.runsMu.Unlock()
	}()

	c.Header("Content-Type", "text/event-stream;charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	_, _ = fmt.Fprint(c.Writer, ": connected\n\n")
	c.Writer.Flush()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-client.Send:
			if !ok {
				return false
			}
			logger.InfoKV(ctx, "sending sse event", "user_id", userID, "event_id", msg.EventID, "type", msg.EventType)
			c.SSEvent("message", notificationResponse{
				ID:        msg.EventID,
				Type:      string(msg.EventType),
				EntityID:  msg.EntityID,
				Metadata:  msg.Metadata,
				CreatedAt: msg.CreatedAt,
			})
			return true
		case <-ticker.C:
			c.SSEvent("ping", "heartbeat")
			return true
		case <-ctx.Done():
			return false
		}
	})
}
