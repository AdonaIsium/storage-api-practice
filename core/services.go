package core

import "context"

type ProvisionService interface {
	RequestCreateVolume(ctx context.Context, req CreateVolumeRequest) (JobID, error)
	RequestExpandVolume(ctx context.Context, id VolumeID, newSize int64, idemKey string) (JobID, error)
	RequestDeleteVolume(ctx context.Context, id VolumeID, idemKey string) (JobID, error)
}

type ConnectivityService interface {
	RequestCreateHost(ctx context.Context, req CreateHostRequest) (JobID, error)
	RequestMapVolume(ctx context.Context, vid VolumeID, hid HostID, opts MapOptions, idemKey string) (JobID, error)
	RequestUnmap(ctx context.Context, mappingID MappingID, idemKey string) (JobID, error)
}

type JobService interface {
	Get(ctx context.Context, id JobID) (*Job, error)
	List(ctx context.Context, state []JobState, limit int) ([]*Job, error)
}
