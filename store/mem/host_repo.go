package mem

import (
	"context"
	"sync"

	"github.com/AdonaIsium/storage-api-practice/core"
	"github.com/AdonaIsium/storage-api-practice/store"
)

type HostRepo struct {
	mu         sync.RWMutex
	byID       map[core.HostID]*core.Host
	byName     map[string]core.HostID
	byIdentity map[string]core.HostID
}

func NewHostRepo() *HostRepo {
	panic("TODO")
}

func (r *HostRepo) Save(ctx context.Context, h *core.Host) error {
	panic("TODO")
}

func (r *HostRepo) Get(ctx context.Context, id core.HostID) (*core.Host, error) {
	panic("TODO")
}

func (r *HostRepo) GetByName(ctx context.Context, name string) (*core.Host, error) {
	panic("TODO")
}

func (r *HostRepo) List(ctx context.Context, f store.HostFilter) ([]*core.Host, error) {
	panic("TODO")
}

func (r *HostRepo) Delete(ctx context.Context, id core.HostID) error {
	panic("TODO")
}

func (r *HostRepo) UpdateCAS(ctx context.Context, h *core.Host, expectedVersion int64) error {
	panic("TODO")
}
