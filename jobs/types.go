package jobs

const (
	JobCreateVolume = "CreateVolume"
	JobExpandVolume = "ExpandVolume"
	JobDeleteVolume = "DeleteVolume"
	JobCreateHost   = "CreateHost"
	JobMapVolume    = "MapVolume"
	JobUnmapVolume  = "UnmapVolume"
)

const (
	StepValidateInput = "ValidateInput"
	StepReserveRecord = "ReserveRecord"
	StepCallDriver    = "CallDriver"
	StepPersistResult = "PersistResult"
	StepEmitAudit     = "EmitAudit"
)
