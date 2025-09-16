package jobs

// Plan returns the ordered step names for a given job type.
func Plan(jobName string) []string {
    switch jobName {
    case JobCreateVolume:
        return []string{StepValidateInput, StepCallDriver, StepPersistResult, StepEmitAudit}
    case JobExpandVolume:
        return []string{StepValidateInput, StepCallDriver, StepPersistResult, StepEmitAudit}
    case JobDeleteVolume:
        return []string{StepValidateInput, StepCallDriver, StepPersistResult, StepEmitAudit}
    case JobCreateHost:
        return []string{StepValidateInput, StepCallDriver, StepPersistResult, StepEmitAudit}
    case JobMapVolume:
        return []string{StepValidateInput, StepCallDriver, StepPersistResult, StepEmitAudit}
    case JobUnmapVolume:
        return []string{StepValidateInput, StepCallDriver, StepPersistResult, StepEmitAudit}
    default:
        return nil
    }
}
