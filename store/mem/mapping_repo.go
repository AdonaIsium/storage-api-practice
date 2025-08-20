package mem

import (
	"context"
	"sync"

	"github.com/AdonaIsium/storage-api-practice/core"
	"github.com/AdonaIsium/storage-api-practice/store"
)

type MappingRepo struct {
	mu     sync.RWMutex
	byID   map[core.MappingID]*core.Mapping
	byVH   map[string]core.MappingID
	byVol  map[core.VolumeID][]core.MappingID
	byHost map[core.HostID][]core.MappingID
}

func NewMappingRepo() *MappingRepo {
	panic("TODO")
}

func (r *MappingRepo) Save(ctx context.Context, m *core.Mapping) error {
	panic("TODO")
}

func (r *MappingRepo) Get(ctx context.Context, id core.MappingID) (*core.Mapping, error) {
	panic("TODO")
}

func (r *MappingRepo) GetByName(ctx context.Context, _ string) (*core.Mapping, error) {
	panic("TODO")
}

func (r *MappingRepo) List(ctx context.Context, f store.MappingFilter) ([]*core.Mapping, error) {
	panic("TODO")
}

func (r *MappingRepo) Delete(ctx context.Context, id core.MappingID) error {
	panic("TODO")
}

func (r *MappingRepo) UpdateCAS(ctx context.Context, _ *core.Mapping, _ int64) error {
	panic("TODO")
}
