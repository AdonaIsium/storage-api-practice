package drivers

import (
	"context"

	"github.com/AdonaIsium/storage-api-practice/core"
)

type Driver interface {
	CreateVolume(ctx context.Context, spec CreateVolumeSpec) (core.Volume, error)
	ExpandVolume(ctx context.Context, id core.VolumeID, newSize int64) error
	DeleteVolume(ctx context.Context, id core.VolumeID) error

	CreateHost(ctx context.Context, spec CreateHostSpec) (core.Host, error)
	MapVolume(ctx context.Context, vol core.VolumeID, host core.HostID, opts MapOpts) (core.Mapping, error)
	UnmapVolume(ctx context.Context, mapping core.MappingID) error

	Health(ctx context.Context) (DriverHealth, error)
}
