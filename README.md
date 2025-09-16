# storage-api-practice

Mock storage orchestration showcasing Redfish/Swordfish-style endpoints, a driver abstraction, async jobs with idempotency, and in-memory stores. Prioritizes Pure-like flows (volumes, hosts/initiators, mappings) via a simple mock.

Quick start
- Run server: `go run ./cmd/swordfish-mock`
- Base URL: `http://localhost:8080`
- Root: `GET /redfish/v1`
- Storage service: `GET /redfish/v1/StorageServices/mock`
- Volumes:
  - List: `GET /redfish/v1/StorageServices/mock/Volumes`
  - Create: `POST /redfish/v1/StorageServices/mock/Volumes` with JSON `{ "Name": "vol1", "CapacityBytes": 1073741824, "ThinProvisioned": true }`
  - Response: `202 Accepted` with `Location: /redfish/v1/TaskService/Tasks/{jobID}`
- Tasks (jobs):
  - `GET /redfish/v1/TaskService/Tasks/{jobID}` returns task state: `Running|Completed|Exception`.

Architecture
- `core/`: domain models + ProvisionService (CreateVolume enqueues job with idempotency).
- `drivers/`: driver interface; `drivers/sim` provides a mock array.
- `jobs/`: job planner, step registry, runner, worker.
- `store/mem`: in-memory repos for volumes/hosts/mappings/jobs/audit/idem.
- `api/swordfish`: minimal Redfish/Swordfish-ish HTTP facade.
- `cmd/swordfish-mock`: wires everything and runs the server.

Notes
- This is a teaching mock, not a full Swordfish implementation. It returns 202 + Task for volume creates to illustrate async job patterns similar to real arrays.
- Extendable to add Hosts/Endpoints and Mappings via similar routes.
