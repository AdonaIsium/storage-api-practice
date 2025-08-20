package store

import (
	"context"

	"github.com/AdonaIsium/storage-api-practice/core"
)

type VolumeRepo interface {
	Save(ctx context.Context, v *core.Volume) error
	Get(ctx context.Context, id core.VolumeID) (*core.Volume, error)
	GetByName(ctx context.Context, name string) (*core.Volume, error)
	List(ctx context.Context, filter VolumeFilter) ([]*core.Volume, error)
	Delete(ctx context.Context, id core.VolumeID) error
	UpdateCAS(ctx context.Context, v *core.Volume, expectedVersion int64) error
}

type HostRepo interface {
	Save(ctx context.Context, v *core.Host) error
	Get(ctx context.Context, id core.HostID) (*core.Host, error)
	GetByName(ctx context.Context, name string) (*core.Host, error)
	List(ctx context.Context, filter HostFilter) ([]*core.Host, error)
	Delete(ctx context.Context, id core.HostID) error
	UpdateCAS(ctx context.Context, v *core.Host, expectedVersion int64)
}

type MappingRepo interface {
	Save(ctx context.Context, v *core.Mapping) error
	Get(ctx context.Context, id core.MappingID) (*core.Mapping, error)
	GetByName(ctx context.Context, name string) (*core.Mapping, error)
	List(ctx context.Context, filter MappingFilter) ([]*core.Mapping, error)
	Delete(ctx context.Context, id core.MappingID) error
	UpdateCAS(ctx context.Context, v *core.Mapping, expectedVersion int64)
}

type JobRepo interface {
	Create(ctx context.Context, j *core.Job) error
	Get(ctx context.Context, id core.JobID) (*core.Job, error)
	Update(ctx context.Context, j *core.Job) error
	List(ctx context.Context, states []core.JobState, limit int) ([]*core.Job, error)
}

type IdemRepo interface {
	Remember(ctx context.Context, op, key string, jobID core.JobID) error
	Recall(ctx context.Context, op, key string) (jobID core.JobID, found bool, err error)
}

type AuditRepo interface {
	Append(ctx context.Context, e *core.AuditEvent) error
	List(ctx context.Context, limit int) ([]*core.AuditEvent, error)
}
