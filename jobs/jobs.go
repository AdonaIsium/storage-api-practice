package jobs

import (
	"context"

	"github.com/AdonaIsium/storage-api-practice/core"
)

type Runner interface {
	Enqueue(ctx context.Context, j *core.Job) (core.JobID, error)
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
