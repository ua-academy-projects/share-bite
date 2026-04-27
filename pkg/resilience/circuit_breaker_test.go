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
		Interval:    10 * time.Millisecond,
		Timeout:     10 * time.Millisecond, // very short timeout for test
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
	})

	opErr := errors.New("op failed")

	// 1st failure
	err := cb.Execute(func() error { return opErr })
	assert.ErrorIs(t, err, opErr)
	assert.Equal(t, gobreaker.StateClosed, cb.State())

	// 2nd failure - should trip
	err = cb.Execute(func() error { return opErr })
	assert.ErrorIs(t, err, opErr)
	assert.Equal(t, gobreaker.StateOpen, cb.State())

	// 3rd request should fail immediately with circuit open
	err = cb.Execute(func() error { return nil })
	assert.ErrorIs(t, err, ErrCircuitOpen)

	// Wait for timeout
	time.Sleep(15 * time.Millisecond)

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
