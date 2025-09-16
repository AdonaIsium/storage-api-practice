package mem

import (
    "context"
    "sync"

    "github.com/AdonaIsium/storage-api-practice/core"
    "github.com/AdonaIsium/storage-api-practice/store"
)

type JobRepo struct {
	mu   sync.RWMutex
	byID map[core.JobID]*core.Job
}

func NewJobRepo() *JobRepo { return &JobRepo{byID: make(map[core.JobID]*core.Job)} }

func (r *JobRepo) Create(ctx context.Context, j *core.Job) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, exists := r.byID[j.ID]; exists {
        return store.ErrConflict
    }
    cp := *j
    r.byID[j.ID] = &cp
    return nil
}

func (r *JobRepo) Get(ctx context.Context, id core.JobID) (*core.Job, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    j, ok := r.byID[id]
    if !ok {
        return nil, store.ErrNotFound
    }
    cp := *j
    return &cp, nil
}

func (r *JobRepo) Update(ctx context.Context, j *core.Job) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.byID[j.ID]; !ok {
        return store.ErrNotFound
    }
    cp := *j
    r.byID[j.ID] = &cp
    return nil
}

func (r *JobRepo) List(ctx context.Context, states []core.JobState, limit int) ([]*core.Job, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    filter := map[core.JobState]struct{}{}
    for _, s := range states {
        filter[s] = struct{}{}
    }
    out := []*core.Job{}
    for _, j := range r.byID {
        if len(filter) > 0 {
            if _, ok := filter[j.State]; !ok {
                continue
            }
        }
        cp := *j
        out = append(out, &cp)
        if limit > 0 && len(out) >= limit {
            break
        }
    }
    return out, nil
}
