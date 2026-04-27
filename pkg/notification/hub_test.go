package notification

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockSubscriber struct {
	mu         sync.Mutex
	msgChans   map[string]chan Message
	readyChans map[string]chan struct{}
}

func newMockSubscriber() *mockSubscriber {
	return &mockSubscriber{
		msgChans:   make(map[string]chan Message),
		readyChans: make(map[string]chan struct{}),
	}
}

func (m *mockSubscriber) Subscribe(ctx context.Context, ch string) (<-chan Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.msgChans[ch]; !ok {
		m.msgChans[ch] = make(chan Message, 10)
		if ready, exists := m.readyChans[ch]; exists {
			close(ready)
		}
	}
	return m.msgChans[ch], nil
}

func (m *mockSubscriber) publish(ch string, msg Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if msgCh, ok := m.msgChans[ch]; ok {
		msgCh <- msg
	}
}

func TestHub_RegisterAndUnregister(t *testing.T) {
	t.Parallel()

	sub := newMockSubscriber()
	hub := NewHub(sub)

	client := &Client{
		UserID: "user1",
		Send:   make(chan Message, 1),
	}

	hub.Register(client)

	hub.mu.RLock()
	_, exists := hub.clients["user1"][client]
	hub.mu.RUnlock()
	assert.True(t, exists)

	hub.Unregister(client)

	hub.mu.RLock()
	_, exists = hub.clients["user1"]
	hub.mu.RUnlock()
	assert.False(t, exists)

	// Sending channel should be closed
	_, ok := <-client.Send
	assert.False(t, ok)
}

func TestHub_Run(t *testing.T) {
	t.Parallel()

	sub := newMockSubscriber()
	hub := NewHub(sub)

	client := &Client{
		UserID: "user1",
		Send:   make(chan Message, 1),
	}
	hub.Register(client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	readyChan := make(chan struct{})
	sub.mu.Lock()
	sub.readyChans["user1"] = readyChan
	sub.mu.Unlock()

	go hub.Run(ctx, "user1")

	select {
	case <-readyChan:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for subscriber")
	}

	msg := Message{
		UserID:    "user1",
		Type:      PostLiked,
		Data:      "123",
		CreatedAt: time.Now(),
	}

	sub.publish("user1", msg)

	select {
	case received := <-client.Send:
		assert.Equal(t, msg.UserID, received.UserID)
		assert.Equal(t, msg.Type, received.Type)
		assert.Equal(t, msg.Data, received.Data)
	case <-time.After(1 * time.Second):
		t.Fatal("did not receive message")
	}
}
