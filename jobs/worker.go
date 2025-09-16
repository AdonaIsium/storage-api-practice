package jobs

import (
    "context"
    "errors"
    "time"

    "github.com/AdonaIsium/storage-api-practice/core"
)

func (r *runner) execJob(ctx context.Context, id core.JobID) {
    j, err := r.deps.Jobs.Get(ctx, id)
    if err != nil {
        return
    }
    j.State = core.JobRunning
    j.UpdatedAt = time.Now()
    _ = r.deps.Jobs.Update(ctx, j)

    steps := Plan(j.Name)
    if len(steps) == 0 {
        j.State = core.JobFailed
        j.ErrorCode = "UNKNOWN_JOB"
        j.ErrorMsg = "no plan for job"
        j.UpdatedAt = time.Now()
        _ = r.deps.Jobs.Update(ctx, j)
        return
    }

    for _, s := range steps {
        dur, err := r.runStep(ctx, j, s)
        _ = dur // currently unused; could log
        if err != nil {
            j.State = core.JobFailed
            j.ErrorCode = "STEP_ERROR"
            j.ErrorMsg = err.Error()
            j.UpdatedAt = time.Now()
            _ = r.deps.Jobs.Update(ctx, j)
            return
        }
    }
    j.State = core.JobSucceeded
    j.UpdatedAt = time.Now()
    _ = r.deps.Jobs.Update(ctx, j)
}

func (r *runner) runStep(ctx context.Context, j *core.Job, stepName string) (dur time.Duration, err error) {
    start := time.Now()
    defer func() { dur = time.Since(start) }()
    f, ok := r.reg.Get(j.Name, stepName)
    if !ok {
        return 0, errors.New("step not registered")
    }
    r.markStep(j, stepName, string(core.StepOK), "started", start, time.Time{})
    if err := f(ctx, r.deps, j); err != nil {
        r.markStep(j, stepName, string(core.StepERR), err.Error(), start, time.Now())
        return 0, err
    }
    r.markStep(j, stepName, string(core.StepOK), "ok", start, time.Now())
    return dur, nil
}

func (r *runner) markStep(j *core.Job, stepName, status, detail string, started, ended time.Time) {
    // update or append
    found := false
    for i := range j.Steps {
        if j.Steps[i].Name == stepName {
            j.Steps[i].Status = core.StepStatus(status)
            if !started.IsZero() {
                j.Steps[i].StartedAt = started
            }
            if !ended.IsZero() {
                j.Steps[i].EndedAt = ended
            }
            j.Steps[i].Detail = detail
            found = true
            break
        }
    }
    if !found {
        j.Steps = append(j.Steps, core.JobStep{
            Name:      stepName,
            Status:    core.StepStatus(status),
            StartedAt: started,
            EndedAt:   ended,
            Detail:    detail,
        })
    }
    j.UpdatedAt = time.Now()
    _ = r.deps.Jobs.Update(context.Background(), j)
}
