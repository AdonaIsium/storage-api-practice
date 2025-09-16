package sim

import (
    "context"
    "math/rand"
    "sync"

    "github.com/AdonaIsium/storage-api-practice/core"
    "github.com/AdonaIsium/storage-api-practice/drivers"
)

type SimDriver struct {
    cfg Config
    mu  sync.Mutex
    rng *rand.Rand
}

func New(cfg Config) *SimDriver {
    // ignore validation error; mock environment
    if cfg.RNGSeed == 0 {
        cfg.RNGSeed = 1
    }
    src := rand.New(rand.NewSource(cfg.RNGSeed))
    return &SimDriver{cfg: cfg, rng: src}
}

func (d *SimDriver) CreateVolume(ctx context.Context, spec drivers.CreateVolumeSpec) (core.Volume, error) {
    if err := d.sleepWithJitter(ctx); err != nil { return core.Volume{}, err }
    if err := d.maybeFail("CreateVolume"); err != nil { return core.Volume{}, err }
    now := d.now()
    return core.Volume{
        ID:        spec.ID,
        Name:      spec.Name,
        SizeBytes: spec.SizeBytes,
        Thin:      spec.Thin,
        QosPolicy: spec.QosPolicy,
        Version:   1,
        CreatedAt: now,
        UpdatedAt: now,
        Tags:      spec.Tags,
    }, nil
}

func (d *SimDriver) ExpandVolume(ctx context.Context, id core.VolumeID, newSize int64) error {
    if err := d.sleepWithJitter(ctx); err != nil { return err }
    return d.maybeFail("ExpandVolume")
}

func (d *SimDriver) DeleteVolume(ctx context.Context, id core.VolumeID) error {
    if err := d.sleepWithJitter(ctx); err != nil { return err }
    return d.maybeFail("DeleteVolume")
}

func (d *SimDriver) CreateHost(ctx context.Context, spec drivers.CreateHostSpec) (core.Host, error) {
    if err := d.sleepWithJitter(ctx); err != nil { return core.Host{}, err }
    if err := d.maybeFail("CreateHost"); err != nil { return core.Host{}, err }
    now := d.now()
    return core.Host{ID: spec.ID, Name: spec.Name, Identities: spec.Identities, CreatedAt: now, UpdatedAt: now}, nil
}

func (d *SimDriver) MapVolume(ctx context.Context, vol core.VolumeID, host core.HostID, opts drivers.MapOpts) (core.Mapping, error) {
    if err := d.sleepWithJitter(ctx); err != nil { return core.Mapping{}, err }
    if err := d.maybeFail("MapVolume"); err != nil { return core.Mapping{}, err }
    lun := 0
    if opts.LUN != nil { lun = *opts.LUN } else { lun = d.chooseLUN() }
    return core.Mapping{ID: core.MappingID(string(vol)+"-"+string(host)), VolumeID: vol, HostID: host, LUN: lun, CreatedAt: d.now()}, nil
}

func (d *SimDriver) UnmapVolume(ctx context.Context, mapping core.MappingID) error {
    if err := d.sleepWithJitter(ctx); err != nil { return err }
    return d.maybeFail("UnmapVolume")
}

func (d *SimDriver) Health(ctx context.Context) (drivers.DriverHealth, error) {
    return drivers.DriverHealth{Ready: true, Detail: "sim"}, nil
}
