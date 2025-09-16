package core

import "context"

type VolumeFilter struct{ NameEquals, TagKey, TagValue string }
type HostFilter struct{ NameEquals, Identity string }
type MappingFilter struct{ VolumeID, HostID string }

type VolumeRepo interface {
	Save(context.Context, *Volume) error
	Get(context.Context, VolumeID) (*Volume, error)
	GetByName(context.Context, string) (*Volume, error)
	List(context.Context, VolumeFilter) ([]*Volume, error)
	Delete(context.Context, VolumeID) error
	UpdateCAS(context.Context, *Volume, int64) error
}
type HostRepo interface {
	Save(context.Context, *Host) error
	Get(context.Context, HostID) (*Host, error)
	GetByName(context.Context, string) (*Host, error)
	List(context.Context, HostFilter) ([]*Host, error)
	Delete(context.Context, HostID) error
	UpdateCAS(context.Context, *Host, int64) error
}
type MappingRepo interface {
	Save(context.Context, *Mapping) error
	Get(context.Context, MappingID) (*Mapping, error)
	GetByName(context.Context, string) (*Mapping, error) // may return not found
	List(context.Context, MappingFilter) ([]*Mapping, error)
	Delete(context.Context, MappingID) error
	UpdateCAS(context.Context, *Mapping, int64) error
}
type JobRepo interface {
	Create(context.Context, *Job) error
	Get(context.Context, JobID) (*Job, error)
	Update(context.Context, *Job) error
	List(context.Context, []JobState, int) ([]*Job, error)
}
type IdemRepo interface {
	Remember(context.Context, string, string, JobID) error // op, key -> jobID
	Recall(context.Context, string, string) (JobID, bool, error)
}
type AuditRepo interface {
	Append(context.Context, *AuditEvent) error
	List(context.Context, int) ([]*AuditEvent, error)
}
