package sim

import (
	"context"
	"time"
)

// sleepWithJitter waits for a random duration between cfg.MinDelay and cfg.MaxDelay,
// but returns early if ctx is done.
func (d *SimDriver) sleepWithJitter(ctx context.Context) error {
	panic("TODO")
}

// maybeFail returns a transient error based on cfg.FailProb; otherwise nil.
// op is a short string like "CreateVolume" for logging/metrics if you add them.
func (d *SimDriver) maybeFail(op string) error {
	panic("TODO")
}

// chooseLUN decides a LUN when opts.LUN == nil (for MapVolume). Keep it simple now
func (d *SimDriver) chooseLUN() int {
	panic("TODO")
}

// now returns a time to stamp resources
func (d *SimDriver) now() time.Time {
	panic("TODO")
}
