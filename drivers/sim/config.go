package sim

import (
	"time"
)

type Config struct {
	MinDelay time.Duration
	MaxDelay time.Duration
	FailProb float64
	RNGSeend int64
}

func (c Config) Validate() error {
	panic("TODO")
}
