package core

import (
	"context"
	"time"

	"github.com/AdonaIsium/storage-api-practice/store"
)

type JobEnqueuer interface {
	Enqueue(ctx context.Context, j *Job) (JobID, error)
}

type ProvisionDeps struct {
	Volumes     store.VolumeRepo
	Jobs        store.JobRepo
	Idempotency store.IdemRepo
	Audit       store.AuditRepo
	Runner      JobEnqueuer
	Now         func() time.Time
}
