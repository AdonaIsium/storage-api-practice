package jobs

import (
    "context"
    "time"

    "github.com/AdonaIsium/storage-api-practice/core"
)

type Runner interface {
	Enqueue(ctx context.Context, j *core.Job) (core.JobID, error)
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

type Config struct {
	Workers      int // number of parallel workers
	PollInterval time.Duration
	// Possibly Add: per-job timeout
}

type runner struct {
    cfg  Config
    deps Deps
    reg  Registry
    q    chan core.JobID
    // done chan struct{} // for Shutdown
}

func New(cfg Config, deps Deps, reg Registry) Runner {
    if cfg.Workers <= 0 { cfg.Workers = 1 }
    if cfg.PollInterval <= 0 { cfg.PollInterval = 100 * time.Millisecond }
    return &runner{cfg: cfg, deps: deps, reg: reg, q: make(chan core.JobID, 128)}
}

// Enqueue persists the job (Jobs.Create) and pushes its ID into q
func (r *runner) Enqueue(ctx context.Context, j *core.Job) (core.JobID, error) {
    // Note: ProvisionService already created Job; but we handle both cases safely
    if j.CreatedAt.IsZero() {
        now := time.Now()
        j.CreatedAt = now
        j.UpdatedAt = now
    }
    // Try create; ignore conflict and proceed
    _ = r.deps.Jobs.Create(ctx, j)
    select {
    case r.q <- j.ID:
    case <-ctx.Done():
        return "", ctx.Err()
    }
    return j.ID, nil
}

// Start launches worker goroutines that read q and execute jobs
func (r *runner) Start(ctx context.Context) error {
    for i := 0; i < r.cfg.Workers; i++ {
        go func() {
            for {
                select {
                case id := <-r.q:
                    r.execJob(ctx, id)
                case <-ctx.Done():
                    return
                }
            }
        }()
    }
    return nil
}

// Shutdown akss workers to stop after current step completes
func (r *runner) Shutdown(ctx context.Context) error {
    // cooperative via context; caller cancels ctx passed to Start
    return nil
}
