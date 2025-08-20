package mem

import (
	"context"
	"sync"

	"github.com/AdonaIsium/storage-api-practice/core"
)

type IdemRepo struct {
	mu   sync.RWMutex
	byOp map[string]map[string]core.JobID
}

func NewIdemRepo() *IdemRepo {
	panic("TODO")
}

func (r *IdemRepo) Remember(ctx context.Context, op, key string, jobID core.JobID) error {
	panic("TODO")
}

func (r *IdemRepo) Recall(ctx context.Context, op, key string) (core.JobID, bool, error) {
	panic("TODO")
}
