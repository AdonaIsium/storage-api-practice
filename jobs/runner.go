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
	panic("TODO")
}

// Enqueue persists the job (Jobs.Create) and pushes its ID into q
func (r *runner) Enqueue(ctx context.Context, j *core.Job) (core.JobID, error) {
	panic("TODO")
}

// Start launches worker goroutines that read q and execute jobs
func (r *runner) Start(ctx context.Context) error {
	panic("TODO")
}

// Shutdown akss workers to stop after current step completes
func (r *runner) Shutdown(ctx context.Context) error {
	panic("TODO")
}
