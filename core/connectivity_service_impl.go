package core

import (
    "context"
    "encoding/json"
    "fmt"
    "strconv"
    "time"

    "github.com/AdonaIsium/storage-api-practice/pkg/ids"
)

type connectivityService struct{ d ConnectivityDeps }

type ConnectivityDeps struct {
    Volumes     VolumeRepo
    Hosts       HostRepo
    Mappings    MappingRepo
    Jobs        JobRepo
    Idempotency IdemRepo
    Audit       AuditRepo
    Runner      JobEnqueuer
    Now         func() time.Time
}

func NewConnectivityService(d ConnectivityDeps) ConnectivityService { return &connectivityService{d: d} }

const (
    opCreateHost  = "CreateHost"
    opMapVolume   = "MapVolume"
    opUnmapVolume = "UnmapVolume"
)

func (s *connectivityService) RequestCreateHost(ctx context.Context, req CreateHostRequest) (JobID, error) {
    if req.Name == "" || len(req.Identities) == 0 {
        return "", fmt.Errorf("name and identities required")
    }
    now := s.d.Now()
    hostID := HostID(ids.NewID())

    // encode identities
    buf, _ := json.Marshal(req.Identities)
    params := map[string]string{
        "host_id":        string(hostID),
        "name":           req.Name,
        "identities_json": string(buf),
    }
    job := &Job{ID: JobID(ids.NewID()), Name: opCreateHost, State: JobPending, Params: params, CorrelationID: ids.NewCorrelationID(), CreatedAt: now, UpdatedAt: now}
    if err := s.d.Jobs.Create(ctx, job); err != nil { return "", err }
    if _, err := s.d.Runner.Enqueue(ctx, job); err != nil { return "", err }
    return job.ID, nil
}

func (s *connectivityService) RequestMapVolume(ctx context.Context, vid VolumeID, hid HostID, opts MapOptions, idemKey string) (JobID, error) {
    if idemKey == "" { return "", fmt.Errorf("idempotency key required") }
    if jid, found, err := s.d.Idempotency.Recall(ctx, opMapVolume, idemKey); err != nil { return "", err } else if found { return jid, nil }
    now := s.d.Now()
    mappingID := MappingID(ids.NewID())
    params := map[string]string{
        "mapping_id": string(mappingID),
        "volume_id":  string(vid),
        "host_id":    string(hid),
    }
    if opts.LUN != nil { params["lun"] = strconv.Itoa(*opts.LUN) }
    job := &Job{ID: JobID(ids.NewID()), Name: opMapVolume, State: JobPending, Params: params, IdempotencyKey: idemKey, CorrelationID: ids.NewCorrelationID(), CreatedAt: now, UpdatedAt: now}
    if err := s.d.Jobs.Create(ctx, job); err != nil { return "", err }
    if err := s.d.Idempotency.Remember(ctx, opMapVolume, idemKey, job.ID); err != nil { return "", err }
    if _, err := s.d.Runner.Enqueue(ctx, job); err != nil { return "", err }
    return job.ID, nil
}

func (s *connectivityService) RequestUnmap(ctx context.Context, mappingID MappingID, idemKey string) (JobID, error) {
    if idemKey == "" { return "", fmt.Errorf("idempotency key required") }
    if jid, found, err := s.d.Idempotency.Recall(ctx, opUnmapVolume, idemKey); err != nil { return "", err } else if found { return jid, nil }
    now := s.d.Now()
    params := map[string]string{ "mapping_id": string(mappingID) }
    job := &Job{ID: JobID(ids.NewID()), Name: opUnmapVolume, State: JobPending, Params: params, IdempotencyKey: idemKey, CorrelationID: ids.NewCorrelationID(), CreatedAt: now, UpdatedAt: now}
    if err := s.d.Jobs.Create(ctx, job); err != nil { return "", err }
    if err := s.d.Idempotency.Remember(ctx, opUnmapVolume, idemKey, job.ID); err != nil { return "", err }
    if _, err := s.d.Runner.Enqueue(ctx, job); err != nil { return "", err }
    return job.ID, nil
}

