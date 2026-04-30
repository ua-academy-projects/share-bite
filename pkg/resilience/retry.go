package resilience

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type RetryConfig struct {
	InitialInterval     time.Duration
	RandomizationFactor float64
	Multiplier          float64
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration
}

type NotifyFn func(err error, nextRetryIn time.Duration)

var ErrNilOperation = errors.New("retry operation cannot be nil")

func NewExponentialBackoff(cfg RetryConfig) *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()

	if cfg.InitialInterval > 0 {
		b.InitialInterval = cfg.InitialInterval
	}
	if cfg.RandomizationFactor > 0 {
		b.RandomizationFactor = cfg.RandomizationFactor
	}
	if cfg.Multiplier > 0 {
		b.Multiplier = cfg.Multiplier
	}
	if cfg.MaxInterval > 0 {
		b.MaxInterval = cfg.MaxInterval
	}
	if cfg.MaxElapsedTime > 0 {
		b.MaxElapsedTime = cfg.MaxElapsedTime
	}

	b.Reset()

	return b
}

func Retry(ctx context.Context, config RetryConfig, operation func() error, notify NotifyFn) error {
	if operation == nil {
		return ErrNilOperation
	}

	b := NewExponentialBackoff(config)
	bo := backoff.BackOff(b)
	if ctx != nil {
		bo = backoff.WithContext(b, ctx)
	}

	if notify != nil {
		err := backoff.RetryNotify(operation, bo, backoff.Notify(notify))
		return unwrapPermanentError(err)
	}

	err := backoff.Retry(operation, bo)
	return unwrapPermanentError(err)
}

func RetryValue[T any](ctx context.Context, config RetryConfig, operation func() (T, error), notify NotifyFn) (T, error) {
	var result T
	if operation == nil {
		return result, ErrNilOperation
	}

	err := Retry(ctx, config, func() error {
		var opErr error
		result, opErr = operation()
		return opErr
	}, notify)

	return result, err
}

func Permanent(err error) error {
	if err == nil {
		return nil
	}

	return backoff.Permanent(err)
}

func IsPermanent(err error) bool {
	if err == nil {
		return false
	}

	var permanentErr *backoff.PermanentError
	return errors.As(err, &permanentErr)
}

func unwrapPermanentError(err error) error {
	if err == nil {
		return nil
	}

	var permanentErr *backoff.PermanentError
	if errors.As(err, &permanentErr) {
		if unwrapped := errors.Unwrap(permanentErr); unwrapped != nil {
			return unwrapped
		}
	}

	return err
}
