<div align="center">

# storage-api-practice

Mock storage orchestration demonstrating Redfish/Swordfish-style endpoints, an abstraction over storage drivers, async jobs with idempotency, and in-memory persistence. It prioritizes Pure-like flows (volumes, hosts/initiators, mappings) and is extendable toward Hitachi semantics.

</div>

## Table of Contents
- [Overview](#overview)
- [Quick Start](#quick-start)
- [API](#api)
- [Architecture](#architecture)
- [End-to-End Flow](#end-to-end-flow)
- [Components](#components)
- [Idempotency](#idempotency)
- [Sim Driver](#sim-driver)
- [How-To: Run & Test](#how-to-run--test)
- [Extending](#extending)

## Overview
This repository models a simplified storage control plane. It exposes a small set of Redfish/Swordfish-ish HTTP endpoints that translate into service calls which enqueue background jobs. Those jobs use a driver abstraction to simulate array behavior and persist results in in-memory repositories.

## Quick Start
```bash
go run ./cmd/swordfish-mock
# server listens on http://localhost:8080
```

## API
Base: `http://localhost:8080`

- Root
  - `GET /redfish/v1` → service links (StorageServices, TaskService)
- Storage service
  - `GET /redfish/v1/StorageServices/mock`
- Volumes
  - `GET /redfish/v1/StorageServices/mock/Volumes` → list
  - `POST /redfish/v1/StorageServices/mock/Volumes` → create (async task)
    ```json
    { "Name": "vol1", "CapacityBytes": 1073741824, "ThinProvisioned": true }
    ```
    Returns `202 Accepted` with `Location: /redfish/v1/TaskService/Tasks/{jobID}`
- Endpoints (Hosts)
  - `GET /redfish/v1/StorageServices/mock/Endpoints` → list
  - `POST /redfish/v1/StorageServices/mock/Endpoints` → create (async task)
    ```json
    { "Name": "host1", "Identifiers": [{ "DurableNameFormat": "iqn", "DurableName": "iqn.1993-08.com.example:01:abc" }] }
    ```
- Mappings (Volume ↔ Endpoint)
  - `GET /redfish/v1/StorageServices/mock/Mappings` → list
  - `POST /redfish/v1/StorageServices/mock/Mappings` → create map (async task)
    ```json
    { "VolumeId": "<vol-id>", "EndpointId": "<endpoint-id>", "LUN": 10 }
    ```
  - `DELETE /redfish/v1/StorageServices/mock/Mappings/{mapping-id}` → unmap (async task)
- Tasks
  - `GET /redfish/v1/TaskService/Tasks/{jobID}` → `Running | Completed | Exception`

## Architecture
- Clean-ish layering with clear ports:
  - `core/` contains domain types and repository/service interfaces.
  - `drivers/` provides a `Driver` interface that abstracts array operations.
  - `jobs/` executes work asynchronously with step pipelines and a worker pool.
  - `store/mem/` offers in-memory repo implementations.
  - `api/swordfish/` is a minimal HTTP facade.

## End-to-End Flow
1. Client sends `POST /Volumes` (or `/Endpoints`, `/Mappings`).
2. HTTP handler builds a service request and (optionally) validates an `Idempotency-Key`.
3. Service enqueues a `Job` (deduplicated via `IdemRepo` for idempotent operations).
4. Workers execute steps for the job: `ValidateInput → CallDriver → PersistResult → EmitAudit`.
5. Results are committed to repos; errors mark the job failed with detail.
6. Client polls `GET /TaskService/Tasks/{id}` for status.

## Components
- API
  - `api/swordfish/server.go` — HTTP handlers: Volumes, Endpoints, Mappings, Tasks.
- Core
  - `core/modeling.go` — Entities (`Volume`, `Host`, `Mapping`, `Job`, `AuditEvent`).
  - `core/dto.go` — Request DTOs.
  - `core/ports.go` — Repo/service interfaces (ports).
  - `core/services.go` — Service interfaces.
  - `core/provision_service_impl.go` — `ProvisionService` (CreateVolume jobs).
  - `core/connectivity_service_impl.go` — `ConnectivityService` (CreateHost/Map/Unmap jobs).
- Jobs
  - `jobs/types.go` — Job and step names.
  - `jobs/planner.go` — Job → ordered steps.
  - `jobs/exec.go` — Step registry + handlers (CreateVolume, CreateHost, MapVolume, UnmapVolume).
  - `jobs/runner.go` — Worker pool and queue.
  - `jobs/worker.go` — Executes steps and updates state.
  - `jobs/deps.go` — Dependencies available to steps.
- Driver
  - `drivers/drivers.go` — Driver interface (Create/Expand/Delete volume; CreateHost; Map/Unmap; Health).
  - `drivers/specs.go` — Driver input specs + health.
  - `drivers/errors.go` — Typed driver errors.
  - `drivers/sim/*` — Deterministic, jittery mock driver.
- Storage (in-memory)
  - `store/mem/*` — Volumes, Hosts, Mappings, Jobs, Idempotency, Audit.
- Wiring
  - `cmd/swordfish-mock/main.go` — Composition root to run the server.
  - `pkg/ids/ids.go` — UUID utilities.

## Idempotency
- Applied to CreateVolume, MapVolume, and UnmapVolume via `IdemRepo`.
- Submissions with the same `(operation, Idempotency-Key)` return the original job ID.
- Mirrors patterns used by Pure/Hitachi for safe retries.

## Sim Driver
- `drivers/sim` uses `MinDelay/MaxDelay` for jitter and `FailProb` for transient failures.
- The driver implements the full `drivers.Driver` interface.
- Configure in `cmd/swordfish-mock/main.go`.

## How-To: Run & Test
See `docs/HOWTO.md` for step-by-step commands using curl and jq.

## Extending
- Add QoS, snapshots, clones; expose them in services and HTTP.
- Add PATCH/DELETE for volumes and hosts; implement expand/delete jobs.
- Add authentication and tags for multi-tenancy.
