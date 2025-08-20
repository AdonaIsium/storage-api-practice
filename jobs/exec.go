package jobs

import (
	"context"

	"github.com/AdonaIsium/storage-api-practice/core"
)

type StepFunc func(ctx context.Context, d Deps, js *core.Job) error

type Registry interface {
	Get(jobName, stepName string) (StepFunc, bool)
}

type InMemoryRegistry struct {
	// map[string]map[string]StepFunc
}

func NewRegistry() *InMemoryRegistry {
	panic("TODO")
}

func (r *InMemoryRegistry) Get(jobName, stepName string) (StepFunc, bool) {
	panic("TODO")
}

func RegisterCreateVolume(r *InMemoryRegistry) {
	panic("TODO")
}

func RegisterExpandVolume(r *InMemoryRegistry) {
	panic("TODO")
}

func RegisterDeleteVolume(r *InMemoryRegistry) {
	panic("TODO")
}

func RegisterCreateHost(r *InMemoryRegistry) {
	panic("TODO")
}

func RegisterMapVolume(r *InMemoryRegistry) {
	panic("TODO")
}

func RegisterUnmapVolume(r *InMemoryRegistry) {
	panic("TODO")
}
