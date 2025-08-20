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
	panic("TODO")
}

func (r *AuditRepo) Append(ctx context.Context, e *core.AuditEvent) error {
	panic("TODO")
}

func (r *AuditRepo) List(ctx context.Context, limit int) ([]*core.AuditEvent, error) {
	panic("TODO")
}
