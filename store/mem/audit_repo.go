package mem

import (
    "context"
    "sync"

    "github.com/AdonaIsium/storage-api-practice/core"
)

type AuditRepo struct {
    mu   sync.RWMutex
    ring []*core.AuditEvent
}

func NewAuditRepo() *AuditRepo {
    return &AuditRepo{ring: make([]*core.AuditEvent, 0, 1024)}
}

func (r *AuditRepo) Append(ctx context.Context, e *core.AuditEvent) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.ring = append(r.ring, e)
    // keep bounded to 1024
    if len(r.ring) > 1024 {
        r.ring = r.ring[len(r.ring)-1024:]
    }
    return nil
}

func (r *AuditRepo) List(ctx context.Context, limit int) ([]*core.AuditEvent, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    n := len(r.ring)
    if limit <= 0 || limit > n {
        limit = n
    }
    start := n - limit
    out := make([]*core.AuditEvent, 0, limit)
    out = append(out, r.ring[start:]...)
    return out, nil
}
