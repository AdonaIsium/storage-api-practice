package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/AdonaIsium/storage-api-practice/api/swordfish"
    "github.com/AdonaIsium/storage-api-practice/core"
    "github.com/AdonaIsium/storage-api-practice/drivers/sim"
    "github.com/AdonaIsium/storage-api-practice/jobs"
    mem "github.com/AdonaIsium/storage-api-practice/store/mem"
)

func main() {
    // deps
    vols := mem.NewVolumeRepo()
    hosts := mem.NewHostRepo()
    maps := mem.NewMappingRepo()
    jobsRepo := mem.NewJobRepo()
    audit := mem.NewAuditRepo()
    idem := mem.NewIdemRepo()

    driver := sim.New(sim.Config{MinDelay: 10 * time.Millisecond, MaxDelay: 50 * time.Millisecond, FailProb: 0})

    deps := jobs.Deps{Driver: driver, Volumes: vols, Hosts: hosts, Mappings: maps, Jobs: jobsRepo, Audit: audit, Idempotency: idem}

    reg := jobs.NewRegistry()
    jobs.RegisterCreateVolume(reg)

    runner := jobs.New(jobs.Config{Workers: 2, PollInterval: 100 * time.Millisecond}, deps, reg)

    prov := core.NewProvisionService(core.ProvisionDeps{
        Volumes:     vols,
        Jobs:        jobsRepo,
        Idempotency: idem,
        Audit:       audit,
        Runner:      runner,
        Now:         time.Now,
    })

    // HTTP
    srv := swordfish.New(deps, prov)
    mux := srv.Routes()

    srvHTTP := &http.Server{Addr: ":8080", Handler: mux}

    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    go func() {
        _ = runner.Start(ctx)
    }()

    go func() {
        log.Printf("listening on %s", srvHTTP.Addr)
        if err := srvHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    <-ctx.Done()
    _ = runner.Shutdown(context.Background())
    _ = srvHTTP.Shutdown(context.Background())
}

