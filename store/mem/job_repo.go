package mem

import (
	"context"
	"sync"

	"github.com/AdonaIsium/storage-api-practice/core"
)

type JobRepo struct {
	mu   sync.RWMutex
	byID map[core.JobID]*core.Job
}

func NewJobRepo() *JobRepo {
	panic("TODO")
}

func (r *JobRepo) Create(ctx context.Context, j *core.Job) error {
	panic("TODO")
}

func (r *JobRepo) Get(ctx context.Context, id core.JobID) (*core.Job, error) {
	panic("TODO")
}

func (r *JobRepo) Update(ctx context.Context, j *core.Job) error {
	panic("TODO")
}

func (r *JobRepo) List(ctx context.Context, states []core.JobState, limit int) ([]*core.Job, error) {
	panic("TODO")
}
