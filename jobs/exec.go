package jobs

import (
    "context"
    "encoding/json"
    "fmt"
    "strconv"
    "time"

    "github.com/AdonaIsium/storage-api-practice/core"
    "github.com/AdonaIsium/storage-api-practice/drivers"
)

type StepFunc func(ctx context.Context, d Deps, js *core.Job) error

type Registry interface {
	Get(jobName, stepName string) (StepFunc, bool)
}

type InMemoryRegistry struct {
    steps map[string]map[string]StepFunc
}

func NewRegistry() *InMemoryRegistry {
    return &InMemoryRegistry{steps: make(map[string]map[string]StepFunc)}
}

func (r *InMemoryRegistry) Get(jobName, stepName string) (StepFunc, bool) {
    if m, ok := r.steps[jobName]; ok {
        f, ok2 := m[stepName]
        return f, ok2
    }
    return nil, false
}

func RegisterCreateVolume(r *InMemoryRegistry) {
    r.must(job(JobCreateVolume,
        step(StepValidateInput, func(ctx context.Context, d Deps, js *core.Job) error {
            name := js.Params["name"]
            sizeStr := js.Params["size_bytes"]
            if name == "" || sizeStr == "" {
                return fmt.Errorf("missing name or size")
            }
            // Ensure no volume with same name exists
            if _, err := d.Volumes.GetByName(ctx, name); err == nil {
                return fmt.Errorf("volume name already exists")
            }
            return nil
        }),
        step(StepCallDriver, func(ctx context.Context, d Deps, js *core.Job) error {
            size, _ := strconv.ParseInt(js.Params["size_bytes"], 10, 64)
            thin, _ := strconv.ParseBool(js.Params["thin"])
            spec := drivers.CreateVolumeSpec{
                ID:        core.VolumeID(js.Params["volume_id"]),
                Name:      js.Params["name"],
                SizeBytes: size,
                Thin:      thin,
                QosPolicy: js.Params["qos_policy"],
                Tags:      nil,
            }
            v, err := d.Driver.CreateVolume(ctx, spec)
            if err != nil {
                return err
            }
            // stash results in job params
            js.Params["created_at"] = v.CreatedAt.Format(time.RFC3339Nano)
            return nil
        }),
        step(StepPersistResult, func(ctx context.Context, d Deps, js *core.Job) error {
            size, _ := strconv.ParseInt(js.Params["size_bytes"], 10, 64)
            thin, _ := strconv.ParseBool(js.Params["thin"])
            now := time.Now()
            v := &core.Volume{
                ID:        core.VolumeID(js.Params["volume_id"]),
                Name:      js.Params["name"],
                SizeBytes: size,
                Thin:      thin,
                QosPolicy: js.Params["qos_policy"],
                Version:   1,
                CreatedAt: now,
                UpdatedAt: now,
                Tags:      nil,
            }
            return d.Volumes.Save(ctx, v)
        }),
        step(StepEmitAudit, func(ctx context.Context, d Deps, js *core.Job) error {
            e := &core.AuditEvent{
                ID:        js.CorrelationID,
                At:        time.Now(),
                Actor:     "system",
                Action:    JobCreateVolume,
                Resources: []string{js.Params["volume_id"]},
                Before:    nil,
                After:     map[string]string{"name": js.Params["name"]},
                Meta:      map[string]string{"job_id": string(js.ID)},
            }
            return d.Audit.Append(ctx, e)
        }),
    ))
}

func RegisterExpandVolume(r *InMemoryRegistry) {
    // For brevity: reuse CallDriver/PersistResult patterns later
}

func RegisterDeleteVolume(r *InMemoryRegistry) {
    // TODO
}

func RegisterCreateHost(r *InMemoryRegistry) {
    r.must(job(JobCreateHost,
        step(StepValidateInput, func(ctx context.Context, d Deps, js *core.Job) error {
            if js.Params["name"] == "" || js.Params["identities_json"] == "" {
                return fmt.Errorf("missing name or identities")
            }
            return nil
        }),
        step(StepCallDriver, func(ctx context.Context, d Deps, js *core.Job) error {
            var ids []core.HostIdentity
            _ = json.Unmarshal([]byte(js.Params["identities_json"]), &ids)
            _, err := d.Driver.CreateHost(ctx, drivers.CreateHostSpec{
                ID:         core.HostID(js.Params["host_id"]),
                Name:       js.Params["name"],
                Identities: ids,
            })
            return err
        }),
        step(StepPersistResult, func(ctx context.Context, d Deps, js *core.Job) error {
            var idsL []core.HostIdentity
            _ = json.Unmarshal([]byte(js.Params["identities_json"]), &idsL)
            now := time.Now()
            h := &core.Host{ID: core.HostID(js.Params["host_id"]), Name: js.Params["name"], Identities: idsL, CreatedAt: now, UpdatedAt: now}
            return d.Hosts.Save(ctx, h)
        }),
        step(StepEmitAudit, func(ctx context.Context, d Deps, js *core.Job) error {
            e := &core.AuditEvent{ID: js.CorrelationID, At: time.Now(), Actor: "system", Action: JobCreateHost, Resources: []string{js.Params["host_id"]}, After: map[string]string{"name": js.Params["name"]}, Meta: map[string]string{"job_id": string(js.ID)}}
            return d.Audit.Append(ctx, e)
        }),
    ))
}

func RegisterMapVolume(r *InMemoryRegistry) {
    r.must(job(JobMapVolume,
        step(StepValidateInput, func(ctx context.Context, d Deps, js *core.Job) error {
            if js.Params["volume_id"] == "" || js.Params["host_id"] == "" {
                return fmt.Errorf("volume_id and host_id required")
            }
            if _, err := d.Volumes.Get(ctx, core.VolumeID(js.Params["volume_id"])); err != nil { return err }
            if _, err := d.Hosts.Get(ctx, core.HostID(js.Params["host_id"])); err != nil { return err }
            return nil
        }),
        step(StepCallDriver, func(ctx context.Context, d Deps, js *core.Job) error {
            var lunPtr *int
            if s := js.Params["lun"]; s != "" {
                if v, err := strconv.Atoi(s); err == nil { lunPtr = &v }
            }
            _, err := d.Driver.MapVolume(ctx, core.VolumeID(js.Params["volume_id"]), core.HostID(js.Params["host_id"]), drivers.MapOpts{LUN: lunPtr})
            return err
        }),
        step(StepPersistResult, func(ctx context.Context, d Deps, js *core.Job) error {
            var lun int
            if s := js.Params["lun"]; s != "" { _ = func() error { v, e := strconv.Atoi(s); if e == nil { lun = v }; return nil }() }
            m := &core.Mapping{ID: core.MappingID(js.Params["mapping_id"]), VolumeID: core.VolumeID(js.Params["volume_id"]), HostID: core.HostID(js.Params["host_id"]), LUN: lun, CreatedAt: time.Now()}
            return d.Mappings.Save(ctx, m)
        }),
        step(StepEmitAudit, func(ctx context.Context, d Deps, js *core.Job) error {
            e := &core.AuditEvent{ID: js.CorrelationID, At: time.Now(), Actor: "system", Action: JobMapVolume, Resources: []string{js.Params["mapping_id"], js.Params["volume_id"], js.Params["host_id"]}, Meta: map[string]string{"job_id": string(js.ID)}}
            return d.Audit.Append(ctx, e)
        }),
    ))
}

func RegisterUnmapVolume(r *InMemoryRegistry) {
    r.must(job(JobUnmapVolume,
        step(StepValidateInput, func(ctx context.Context, d Deps, js *core.Job) error {
            if js.Params["mapping_id"] == "" { return fmt.Errorf("mapping_id required") }
            return nil
        }),
        step(StepCallDriver, func(ctx context.Context, d Deps, js *core.Job) error {
            return d.Driver.UnmapVolume(ctx, core.MappingID(js.Params["mapping_id"]))
        }),
        step(StepPersistResult, func(ctx context.Context, d Deps, js *core.Job) error {
            return d.Mappings.Delete(ctx, core.MappingID(js.Params["mapping_id"]))
        }),
        step(StepEmitAudit, func(ctx context.Context, d Deps, js *core.Job) error {
            e := &core.AuditEvent{ID: js.CorrelationID, At: time.Now(), Actor: "system", Action: JobUnmapVolume, Resources: []string{js.Params["mapping_id"]}, Meta: map[string]string{"job_id": string(js.ID)}}
            return d.Audit.Append(ctx, e)
        }),
    ))
}

// helpers to register steps
type regJob struct{
    name string
    steps map[string]StepFunc
}
func job(name string, s ...regStep) regJob {
    m := make(map[string]StepFunc)
    for _, st := range s { m[st.name] = st.fn }
    return regJob{name: name, steps: m}
}
type regStep struct{ name string; fn StepFunc }
func step(name string, fn StepFunc) regStep { return regStep{name: name, fn: fn} }
func (r *InMemoryRegistry) must(j regJob) {
    if r.steps[j.name] == nil { r.steps[j.name] = map[string]StepFunc{} }
    for n, f := range j.steps { r.steps[j.name][n] = f }
}
