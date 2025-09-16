package core

import "time"

// Types here model our SAN for mocking
type VolumeID string
type HostID string
type MappingID string
type JobID string

type Volume struct {
	ID        VolumeID
	Name      string
	SizeBytes int64
	Thin      bool
	QosPolicy string
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      map[string]string
}

type IdentityType string

const (
	IdentityIQN  IdentityType = "iqn"
	IdentityWWPN IdentityType = "wwpn"
	IdentityNQN  IdentityType = "nqn"
)

type HostIdentity struct {
	Type  IdentityType
	Value string
}

type Host struct {
	ID         HostID
	Name       string
	Identities []HostIdentity
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Mapping struct {
	ID        MappingID
	VolumeID  VolumeID
	HostID    HostID
	LUN       int
	CreatedAt time.Time
}

type JobState string

const (
	JobPending   JobState = "PENDING"
	JobRunning   JobState = "RUNNING"
	JobSucceeded JobState = "SUCCEEDED"
	JobFailed    JobState = "FAILED"
)

type Job struct {
	ID             JobID
	Name           string
	State          JobState
	Steps          []JobStep
	IdempotencyKey string
	CorrelationID  string
	ErrorCode      string
	ErrorMsg       string
	ErrorDetails   map[string]any
	Params         map[string]string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type StepStatus string

const (
	StepOK   StepStatus = "OK"
	StepERR  StepStatus = "ERR"
	StepSKIP StepStatus = "SKIP"
)

type JobStep struct {
	Name      string
	Status    StepStatus
	StartedAt time.Time
	EndedAt   time.Time
	Detail    string
}

type AuditEvent struct {
	ID        string
	At        time.Time
	Actor     string
	Action    string
	Resources []string
	Before    any
	After     any
	Meta      map[string]string
}
