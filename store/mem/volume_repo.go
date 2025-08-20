package mem

import (
	"context"
	"sync"

	"github.com/AdonaIsium/storage-api-practice/core"
	"github.com/AdonaIsium/storage-api-practice/store"
)

type VolumeRepo struct {
	mu     sync.RWMutex
	byID   map[core.VolumeID]*core.Volume
	byName map[string]core.VolumeID
}

func NewVolumeRepo() *VolumeRepo {
	panic("TODO")
}

func (r *VolumeRepo) Save(ctx context.Context, v *core.Volume) error {
	panic("TODO")
}

func (r *VolumeRepo) GetByName(ctx context.Context, name string) (*core.Volume, error) {
	panic("TODO")
}

func (r *VolumeRepo) List(ctx context.Context, f store.VolumeFilter) ([]*core.Volume, error) {
	panic("TODO")
}

func (r *VolumeRepo) Delete(ctx context.Context, id core.VolumeID) error {
	panic("TODO")
}

func (r *VolumeRepo) UpdateCAS(ctx context.Context, v *core.Volume, expectedVersion int64) error {
	panic("TODO")
}
