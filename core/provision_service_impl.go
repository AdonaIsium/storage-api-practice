package core

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/AdonaIsium/storage-api-practice/pkg/ids"
)

type JobEnqueuer interface {
	Enqueue(ctx context.Context, j *Job) (JobID, error)
}

type ProvisionDeps struct {
	Volumes     VolumeRepo
	Jobs        JobRepo
	Idempotency IdemRepo
	Audit       AuditRepo
	Runner      JobEnqueuer
	Now         func() time.Time
}

type provisionService struct {
	d ProvisionDeps
}

func NewProvisionService(d ProvisionDeps) ProvisionService {
	return &provisionService{d: d}
}

var _ ProvisionService = (*provisionService)(nil)

const opCreateVolume = "CreateVolume"

func (s *provisionService) RequestCreateVolume(ctx context.Context, req CreateVolumeRequest, idemKey string) (JobID, error) {
	// require idempotency key
	if idemKey == "" {
		return "", fmt.Errorf("idempotency key required")
	}

	// dedupe by (operation, key)
	if jid, found, err := s.d.Idempotency.Recall(ctx, opCreateVolume, idemKey); err != nil {
		return "", err
	} else if found {
		return jid, nil
	}

	if req.Name == "" {
		return "", fmt.Errorf("name is required")
	}
	if req.SizeBytes <= 0 {
		return "", fmt.Errorf("size_bytes must be > 0")
	}

	now := s.d.Now()
	volID := VolumeID(ids.NewID())

	params := map[string]string{
		"volume_id":  string(volID),
		"name":       req.Name,
		"size_bytes": strconv.FormatInt(req.SizeBytes, 10),
		"thin":       strconv.FormatBool(req.Thin),
		"qos_policy": req.QosPolicy,
		// NOTE: tags not included here, can JSON-encode Params later if needed
	}

	job := &Job{
		ID:             JobID(ids.NewID()),
		Name:           opCreateVolume,
		State:          JobPending,
		Params:         params,
		IdempotencyKey: idemKey,
		CorrelationID:  ids.NewCorrelationID(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Possible future place to tighten things up.
	// Right now if Jobs.Create is successful but Idempotency.Remember is not,
	// We could end up with multiple jobs being created
	// Eventually probably want to wrap both in a transaction or implement
	// some sort of "idempotency lock"
	if err := s.d.Jobs.Create(ctx, job); err != nil {
		return "", err
	}

	if err := s.d.Idempotency.Remember(ctx, opCreateVolume, idemKey, job.ID); err != nil {
		return "", err
	}

	if _, err := s.d.Runner.Enqueue(ctx, job); err != nil {
		return "", err
	}

	return job.ID, nil
}

func (s *provisionService) RequestExpandVolume(ctx context.Context, id VolumeID, newSize int64, idemKey string) (JobID, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *provisionService) RequestDeleteVolume(ctx context.Context, id VolumeID, idemKey string) (JobID, error) {
	return "", fmt.Errorf("not implemented")
}
