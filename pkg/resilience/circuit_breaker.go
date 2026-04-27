package resilience

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

const (
	defaultBreakerMinRequests = 10
	defaultBreakerFailureRate = 0.5
)

var (
	ErrCircuitOpen        = gobreaker.ErrOpenState
	ErrCircuitHalfOpenMax = gobreaker.ErrTooManyRequests
)

type CircuitBreakerConfig struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts gobreaker.Counts) bool
	IsSuccessful  func(err error) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

type CircuitBreaker struct {
	breaker *gobreaker.CircuitBreaker
}

func NewCircuitBreaker(config CircuitBreakerConfig) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:          config.Name,
		MaxRequests:   config.MaxRequests,
		Interval:      config.Interval,
		Timeout:       config.Timeout,
		ReadyToTrip:   config.ReadyToTrip,
		IsSuccessful:  config.IsSuccessful,
		OnStateChange: config.OnStateChange,
	}

	if settings.ReadyToTrip == nil {
		settings.ReadyToTrip = defaultReadyToTrip
	}
	if settings.IsSuccessful == nil {
		settings.IsSuccessful = defaultIsSuccessful
	}

	return &CircuitBreaker{breaker: gobreaker.NewCircuitBreaker(settings)}
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
	if operation == nil {
		return ErrNilOperation
	}

	if cb == nil || cb.breaker == nil {
		return operation()
	}

	_, err := cb.breaker.Execute(func() (interface{}, error) {
		return nil, operation()
	})

	return err
}

func ExecuteValueWithBreaker[T any](cb *CircuitBreaker, operation func() (T, error)) (T, error) {
	var zero T
	if operation == nil {
		return zero, ErrNilOperation
	}

	if cb == nil || cb.breaker == nil {
		return operation()
	}

	result, err := cb.breaker.Execute(func() (interface{}, error) {
		return operation()
	})
	if err != nil {
		return zero, err
	}

	typed, ok := result.(T)
	if !ok {
		return zero, fmt.Errorf("unexpected circuit breaker result type")
	}

	return typed, nil
}

func (cb *CircuitBreaker) State() gobreaker.State {
	if cb == nil || cb.breaker == nil {
		return gobreaker.StateClosed
	}

	return cb.breaker.State()
}

func (cb *CircuitBreaker) Counts() gobreaker.Counts {
	if cb == nil || cb.breaker == nil {
		return gobreaker.Counts{}
	}

	return cb.breaker.Counts()
}

func defaultReadyToTrip(counts gobreaker.Counts) bool {
	if counts.Requests < defaultBreakerMinRequests {
		return false
	}

	failureRate := float64(counts.TotalFailures) / float64(counts.Requests)
	return failureRate >= defaultBreakerFailureRate
}

// defaultIsSuccessful treats any non-nil error as a failure.
// This means context.Canceled and context.DeadlineExceeded will count
// toward failure rates in defaultReadyToTrip.
// Advise callers to provide a custom IsSuccessful when they want to treat
// context cancellations or timeouts as non-failures.
func defaultIsSuccessful(err error) bool {
	return err == nil
}
