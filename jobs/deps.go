package jobs

import (
	"github.com/AdonaIsium/storage-api-practice/drivers"
	"github.com/AdonaIsium/storage-api-practice/store"
)

type Deps struct {
	Driver      drivers.Driver
	Volumes     store.VolumeRepo
	Hosts       store.HostRepo
	Mappings    store.MappingRepo
	Jobs        store.JobRepo
	Audit       store.AuditRepo
	Idempotency store.IdemRepo
}
