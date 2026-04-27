package resilience

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry_SuccessOnFirstTry(t *testing.T) {
	t.Parallel()

	called := 0
	err := Retry(context.Background(), RetryConfig{}, func() error {
		called++
		return nil
	}, nil)

	require.NoError(t, err)
	assert.Equal(t, 1, called)
}

func TestRetry_SuccessAfterRetries(t *testing.T) {
	t.Parallel()

	called := 0
	opErr := errors.New("temp error")

	err := Retry(context.Background(), RetryConfig{
		InitialInterval: 1 * time.Millisecond,
		MaxElapsedTime:  50 * time.Millisecond,
	}, func() error {
		called++
		if called < 3 {
			return opErr
		}
		return nil
	}, nil)

	require.NoError(t, err)
	assert.Equal(t, 3, called)
}

func TestRetry_PermanentError(t *testing.T) {
	t.Parallel()

	called := 0
	opErr := errors.New("fatal error")

	err := Retry(context.Background(), RetryConfig{
		InitialInterval: 1 * time.Millisecond,
		MaxElapsedTime:  50 * time.Millisecond,
	}, func() error {
		called++
		return Permanent(opErr)
	}, nil)

	require.ErrorIs(t, err, opErr)
	assert.Equal(t, 1, called) // Should not retry
}

func TestRetryValue_Success(t *testing.T) {
	t.Parallel()

	called := 0
	val, err := RetryValue(context.Background(), RetryConfig{
		InitialInterval: 1 * time.Millisecond,
	}, func() (string, error) {
		called++
		if called < 2 {
			return "", errors.New("temp error")
		}
		return "success", nil
	}, nil)

	require.NoError(t, err)
	assert.Equal(t, "success", val)
	assert.Equal(t, 2, called)
}
