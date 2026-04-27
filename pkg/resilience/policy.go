package resilience

import (
	"context"
	"errors"
)

type Policy struct {
	RetryConfig  RetryConfig
	RetryNotify  NotifyFn
	Breaker      *CircuitBreaker
	DisableRetry bool
}

func (p Policy) Execute(ctx context.Context, operation func() error) error {
	if operation == nil {
		return ErrNilOperation
	}

	wrapped := operation
	if p.Breaker != nil {
		wrapped = func() error {
			err := p.Breaker.Execute(operation)
			if errors.Is(err, ErrCircuitOpen) || errors.Is(err, ErrCircuitHalfOpenMax) {
				return Permanent(err)
			}

			return err
		}
	}

	if p.DisableRetry {
		return wrapped()
	}

	return Retry(ctx, p.RetryConfig, wrapped, p.RetryNotify)
}

func ExecuteValue[T any](ctx context.Context, policy Policy, operation func() (T, error)) (T, error) {
	if operation == nil {
		var zero T
		return zero, ErrNilOperation
	}

	wrapped := operation
	if policy.Breaker != nil {
		wrapped = func() (T, error) {
			result, err := ExecuteValueWithBreaker(policy.Breaker, operation)
			if errors.Is(err, ErrCircuitOpen) || errors.Is(err, ErrCircuitHalfOpenMax) {
				return result, Permanent(err)
			}

			return result, err
		}
	}

	if policy.DisableRetry {
		return wrapped()
	}

	return RetryValue(ctx, policy.RetryConfig, wrapped, policy.RetryNotify)
}
