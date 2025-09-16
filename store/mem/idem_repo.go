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
    return &IdemRepo{byOp: make(map[string]map[string]core.JobID)}
}

func (r *IdemRepo) Remember(ctx context.Context, op, key string, jobID core.JobID) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    m, ok := r.byOp[op]
    if !ok {
        m = make(map[string]core.JobID)
        r.byOp[op] = m
    }
    if _, exists := m[key]; !exists {
        m[key] = jobID
    }
    return nil
}

func (r *IdemRepo) Recall(ctx context.Context, op, key string) (core.JobID, bool, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    if m, ok := r.byOp[op]; ok {
        if id, ok2 := m[key]; ok2 {
            return id, true, nil
        }
    }
    return "", false, nil
}
