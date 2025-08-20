package jobs

import (
	"context"
	"time"

	"github.com/AdonaIsium/storage-api-practice/core"
)

func (r *runner) execJob(ctx context.Context, id core.JobID) {
	panic("TODO")
}

func (r *runner) runStep(ctx context.Context, j *core.Job, stepName string) (dur time.Duration, err error) {
	panic("TODO")
}

func (r *runner) markStep(j *core.Job, stepName, status, detail string, started, ended time.Time) {
	panic("TODO")
}
