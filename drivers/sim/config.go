package sim

import (
    "fmt"
    "time"
)

type Config struct {
    MinDelay time.Duration
    MaxDelay time.Duration
    FailProb float64
    RNGSeed  int64
}

func (c Config) Validate() error {
    if c.MinDelay < 0 || c.MaxDelay < 0 {
        return fmt.Errorf("delays must be >= 0")
    }
    if c.MaxDelay < c.MinDelay {
        return fmt.Errorf("max delay < min delay")
    }
    if c.FailProb < 0 || c.FailProb > 1 {
        return fmt.Errorf("fail prob must be in [0,1]")
    }
    return nil
}
