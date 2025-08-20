package sim

import (
	"context"
	"math/rand/v2"
	"sync"

	"github.com/AdonaIsium/storage-api-practice/core"
	"github.com/AdonaIsium/storage-api-practice/drivers"
)

type SimDriver struct {
	cfg Config
	mu  sync.Mutex
	rng *rand.Rand
}

func New(cfg Config) *SimDriver {
	panic("TODO")
}

func (d *SimDriver) CreateVolume(ctx context.Context, spec drivers.CreateVolumeSpec) (core.Volume, error) {
	panic("TODO")
}

func (d *SimDriver) ExpandVolume(ctx context.Context, id core.VolumeID, newSize int64) error {
	panic("TODO")
}

func (d *SimDriver) CreateHost(ctx context.Context, spec drivers.CreateHostSpec) (core.Host, error) {
	panic("TODO")
}

func (d *SimDriver) MapVolume(ctx context.Context, vol core.VolumeID, host core.HostID, opts drivers.MapOpts) (core.Mapping, error) {
	panic("TODO")
}

func (d *SimDriver) UnmapVolume(ctx context.Context, mapping core.MappingID) error {
	panic("TODO")
}

func (d *SimDriver) Health(ctx context.Context) (drivers.DriverHealth, error) {
	panic("TODO")
}
