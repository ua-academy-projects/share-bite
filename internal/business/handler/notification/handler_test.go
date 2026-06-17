package notification

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ua-academy-projects/share-bite/internal/middleware"
	"github.com/ua-academy-projects/share-bite/pkg/notification"
)

type sseTestResponseWriter struct {
	*httptest.ResponseRecorder
	closeNotifyChan chan bool
}

func (w *sseTestResponseWriter) CloseNotify() <-chan bool {
	return w.closeNotifyChan
}

type mockSubscriber struct {
	mu         sync.Mutex
	msgChans   map[string]chan notification.Message
	readyChans map[string]chan struct{}
}

func newMockSubscriber() *mockSubscriber {
	return &mockSubscriber{
		msgChans:   make(map[string]chan notification.Message),
		readyChans: make(map[string]chan struct{}),
	}
}

func (m *mockSubscriber) Subscribe(_ context.Context, ch string) (<-chan notification.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.msgChans[ch]; !ok {
		m.msgChans[ch] = make(chan notification.Message, 10)
		if ready, exists := m.readyChans[ch]; exists {
			close(ready)
		}
	}
	return m.msgChans[ch], nil
}

func (m *mockSubscriber) publish(ch string, msg notification.Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if msgCh, ok := m.msgChans[ch]; ok {
		msgCh <- msg
	}
}

func TestStream_SuccessAndMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sub := newMockSubscriber()
	hub := notification.NewHub(sub)

	h := &handler{
		hub:        hub,
		activeRuns: make(map[string]bool),
	}

	r := gin.New()
	const testUserID = "user-uuid-1234"
	r.GET("/stream", func(c *gin.Context) {
		c.Set(middleware.CtxUserID, testUserID)
		h.stream(c)
	})

	readyChan := make(chan struct{})
	sub.mu.Lock()
	sub.readyChans[testUserID] = readyChan
	sub.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/stream", nil)
	recorder := httptest.NewRecorder()
	w := &sseTestResponseWriter{
		ResponseRecorder: recorder,
		closeNotifyChan:  make(chan bool, 1),
	}

	go func() {
		r.ServeHTTP(w, req)
	}()

	select {
	case <-readyChan:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for handler to start hub subscription")
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
	if contentType := w.Header().Get("Content-Type"); !strings.Contains(contentType, "text/event-stream") {
		t.Errorf("expected content type text/event-stream, got %s", contentType)
	}

	msg := notification.NewMessageWithMetadata(
		notification.PostLiked,
		testUserID,
		"actor-uuid-5678",
		"post",
		"entity-uuid-9999",
		map[string]any{"message": "test-payload"},
		time.Now(),
	)

	sub.publish(testUserID, msg)
	time.Sleep(50 * time.Millisecond)
	cancel()
	w.closeNotifyChan <- true
	time.Sleep(10 * time.Millisecond)

	body := w.Body.String()

	if !strings.Contains(body, ": connected") {
		t.Errorf("expected body to contain connection greeting")
	}
	if !strings.Contains(body, "event:post_liked") {
		t.Errorf("expected body to contain SSE event type 'post_liked', got: %s", body)
	}
	if !strings.Contains(body, `"message":"test-payload"`) {
		t.Errorf("expected body to contain JSON encoded metadata 'test-payload', got: %s", body)
	}
}

func TestStream_MissingUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	sub := newMockSubscriber()
	h := &handler{
		hub:        notification.NewHub(sub),
		activeRuns: make(map[string]bool),
	}

	r := gin.New()
	r.GET("/stream", h.stream)

	req, _ := http.NewRequest(http.MethodGet, "/stream", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "user id not found in context") {
		t.Errorf("unexpected body error message: %s", w.Body.String())
	}
}
