package resilience

import (
	"errors"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		Name:        "test-cb",
		MaxRequests: 1,
		Interval:    time.Minute,
		Timeout:     time.Minute,
	})

	called := false
	err := cb.Execute(func() error {
		called = true
		return nil
	})

	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, gobreaker.StateClosed, cb.State())
}

func TestCircuitBreaker_Execute_TripsAndRecovers(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		Name:        "test-cb",
		MaxRequests: 1,
		Interval:    time.Minute,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
	})

	opErr := errors.New("op failed")

	err := cb.Execute(func() error { return opErr })
	assert.ErrorIs(t, err, opErr)
	assert.Equal(t, gobreaker.StateClosed, cb.State())

	err = cb.Execute(func() error { return opErr })
	assert.ErrorIs(t, err, opErr)
	assert.Equal(t, gobreaker.StateOpen, cb.State())

	err = cb.Execute(func() error { return nil })
	assert.ErrorIs(t, err, ErrCircuitOpen)

	time.Sleep(75 * time.Millisecond)

	// Next request is half-open, success should close it
	err = cb.Execute(func() error { return nil })
	require.NoError(t, err)
	assert.Equal(t, gobreaker.StateClosed, cb.State())
}

func TestExecuteValueWithBreaker_Success(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		Name: "test-cb",
	})

	val, err := ExecuteValueWithBreaker(cb, func() (string, error) {
		return "success", nil
	})

	require.NoError(t, err)
	assert.Equal(t, "success", val)
}
