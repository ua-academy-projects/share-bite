package env

import (
	"os"
	"strings"
	"sync"

	caarlos0 "github.com/caarlos0/env/v11"
)

var (
	mu        sync.RWMutex
	globalEnv map[string]string
)

type Options = caarlos0.Options

// Init initializes the global environment with secrets and system environment variables.
func Init(secrets map[string]string) {
	mu.Lock()
	defer mu.Unlock()

	if globalEnv == nil {
		globalEnv = make(map[string]string)
		for _, e := range os.Environ() {
			if key, value, ok := strings.Cut(e, "="); ok {
				globalEnv[key] = value
			}
		}
	}

	for k, v := range secrets {
		globalEnv[k] = v
	}
}

// GetSecret returns a secret by key from the global environment.
func GetSecret(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if globalEnv != nil {
		if v, ok := globalEnv[key]; ok {
			return v
		}
	}
	return os.Getenv(key)
}

// Parse parses environment variables into the provided struct using global environment and optional overrides.
func Parse(v any, opts ...caarlos0.Options) error {
	var opt caarlos0.Options
	if len(opts) > 0 {
		opt = opts[0]
	}

	mu.RLock()
	defer mu.RUnlock()

	if globalEnv != nil {
		mergedEnv := make(map[string]string, len(globalEnv)+len(opt.Environment))

		for k, val := range globalEnv {
			mergedEnv[k] = val
		}

		for k, val := range opt.Environment {
			mergedEnv[k] = val
		}

		opt.Environment = mergedEnv
	}

	return caarlos0.ParseWithOptions(v, opt)
}
