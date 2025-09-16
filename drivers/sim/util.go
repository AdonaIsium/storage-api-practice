package sim

import (
    "context"
    "time"

    "github.com/AdonaIsium/storage-api-practice/drivers"
)

// sleepWithJitter waits for a random duration between cfg.MinDelay and cfg.MaxDelay,
// but returns early if ctx is done.
func (d *SimDriver) sleepWithJitter(ctx context.Context) error {
    min := d.cfg.MinDelay
    max := d.cfg.MaxDelay
    if max < min { max = min }
    dur := min
    if max > min {
        diff := max - min
        // pick [0, diff]
        var n64 int64
        d.mu.Lock()
        if diff == time.Duration(^uint64(0)>>1) { // extremely large; avoid +1 overflow
            n64 = d.rng.Int63()
        } else {
            n64 = d.rng.Int63n(int64(diff) + 1)
        }
        d.mu.Unlock()
        dur = min + time.Duration(n64)
    }
    t := time.NewTimer(dur)
    defer t.Stop()
    select {
    case <-t.C:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// maybeFail returns a transient error based on cfg.FailProb; otherwise nil.
// op is a short string like "CreateVolume" for logging/metrics if you add them.
func (d *SimDriver) maybeFail(op string) error {
    d.mu.Lock()
    p := d.cfg.FailProb
    r := d.rng.Float64()
    d.mu.Unlock()
    if r < p {
        return drivers.NewError(drivers.CodeBusy, op+" transient error", true)
    }
    return nil
}

// chooseLUN decides a LUN when opts.LUN == nil (for MapVolume). Keep it simple now
func (d *SimDriver) chooseLUN() int {
    d.mu.Lock()
    n := d.rng.Intn(256)
    d.mu.Unlock()
    return n
}

// now returns a time to stamp resources
func (d *SimDriver) now() time.Time {
    return time.Now()
}
