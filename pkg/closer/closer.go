package closer

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ua-academy-projects/share-bite/pkg/logger"
)

const (
	defaultShutdownTimeout = time.Minute * 5
)

var instance = newCloser(syscall.SIGTERM, syscall.SIGINT)

type closer struct {
	mu              sync.Mutex
	once            sync.Once
	done            chan struct{}
	funcs           []func(context.Context) error
	shutdownTimeout time.Duration
}

func SetShutdownTimeout(t time.Duration) {
	instance.shutdownTimeout = t
}

func newCloser(sigs ...os.Signal) *closer {
	c := &closer{
		done:            make(chan struct{}),
		shutdownTimeout: defaultShutdownTimeout,
	}
	if len(sigs) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, sigs...)
			defer signal.Stop(ch)
			<-ch

			logger.Logger().Info("starting graceful shutdown shenanigans")

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), c.shutdownTimeout)
			defer shutdownCancel()

			if err := c.closeAll(shutdownCtx); err != nil {
				logger.Logger().Errorf("close all finished with errors: %v", err)
			}
		}()
	}

	return c
}

func Wait() {
	instance.wait()
}

func Add(f ...func(context.Context) error) {
	instance.add(f...)
}

func (c *closer) add(f ...func(context.Context) error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *closer) wait() {
	<-c.done
}

func (c *closer) closeAll(ctx context.Context) error {
	var res error
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			logger.Logger().Info("no funcs to close in closer")
			return
		}

		errs := make(chan error, len(funcs))
		wg := sync.WaitGroup{}

		for _, f := range funcs {
			wg.Add(1)
			go func(f func(context.Context) error) {
				defer wg.Done()

				defer func() {
					if r := recover(); r != nil {
						errs <- errors.New("panic recovered in closer")
					}
				}()

				errs <- f(ctx)
			}(f)
		}

		go func() {
			wg.Wait()
			close(errs)
		}()

		for {
			select {
			case <-ctx.Done():
				logger.Logger().Warn("shutdown timed out")
				res = errors.Join(res, ctx.Err())
				return
			case err, ok := <-errs:
				if !ok {
					if res == nil {
						logger.Logger().Info("graceful shutdown completed successfully")
					}
					return
				}

				if err != nil {
					logger.Logger().Errorf("closer graceful shutdown: %s", err)
					res = errors.Join(res, err)
				}
			}
		}
	})

	return res
}
