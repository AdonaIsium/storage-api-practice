package drivers

import "github.com/AdonaIsium/storage-api-practice/core"

type CreateVolumeSpec struct {
	ID        core.VolumeID
	Name      string
	SizeBytes int64
	Thin      bool
	QosPolicy string
	Tags      map[string]string
}

type CreateHostSpec struct {
	ID         core.HostID
	Name       string
	Identities []core.HostIdentity
}

type MapOpts struct {
	LUN *int
}

type DriverHealth struct {
	Ready  bool
	Detail string
}
