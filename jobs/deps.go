package jobs

import (
	"github.com/AdonaIsium/storage-api-practice/core"
	"github.com/AdonaIsium/storage-api-practice/drivers"
)

type Deps struct {
	Driver      drivers.Driver
	Volumes     core.VolumeRepo
	Hosts       core.HostRepo
	Mappings    core.MappingRepo
	Jobs        core.JobRepo
	Audit       core.AuditRepo
	Idempotency core.IdemRepo
}
