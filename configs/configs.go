package configs

import (
	"sync"

	"github.com/caarlos0/env/v7"
)

var once sync.Once

// NewConfig initializes a new Config object by parsing environment variables.
// The function uses sync.Once to ensure that the initialization happens only once.
// The returned Config object can be used to access the parsed configuration values.
func NewConfig() (cfg Config, err error) {

	once.Do(func() {
		cfg = Config{}
		err = env.Parse(&cfg)
	})

	return
}
