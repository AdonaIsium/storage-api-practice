# How To Run and Test the Mock Storage API

This guide shows how to run the Swordfish-style mock, create volumes and endpoints (hosts), map/unmap them, and verify tasks.

## Prerequisites
- Go 1.21+ (ideally 1.22)
- curl, jq (optional, for pretty output)

## Start the Server

```bash
# from the repo root
go run ./cmd/swordfish-mock
# server listens on http://localhost:8080
```

## Quick Smoke Test

1) Root endpoint
```bash
curl -s http://localhost:8080/redfish/v1 | jq
```

2) Create a volume (async task)
```bash
curl -i -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Volumes \
  -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: 00000000-0000-4000-8000-000000000001' \
  -d '{"Name":"vol1","CapacityBytes":1073741824,"ThinProvisioned":true}'
```
- Copy the `Location` header (e.g., `/redfish/v1/TaskService/Tasks/<jobID>`)

3) Poll the task until `Completed`
```bash
curl -s http://localhost:8080/redfish/v1/TaskService/Tasks/<jobID> | jq
```

4) Verify volume exists
```bash
curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Volumes | jq
```

## Endpoints (Hosts)
1) Create an endpoint (host)
```bash
curl -i -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Endpoints \
  -H 'Content-Type: application/json' \
  -d '{"Name":"host1","Identifiers":[{"DurableNameFormat":"iqn","DurableName":"iqn.1993-08.com.example:01:abc"}]}'
```
- Poll task to `Completed`, then list:
```bash
curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Endpoints | jq
```

## Mappings
1) Capture IDs:
```bash
VOL=$(curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Volumes | jq -r '.Members[0].Id')
EP=$(curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Endpoints | jq -r '.Members[0].Id')
echo "$VOL $EP"
```

2) Map volume to endpoint (async)
```bash
curl -i -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Mappings \
  -H 'Content-Type: application/json' \
  -H 'Idempotency-Key: 00000000-0000-4000-8000-000000000002' \
  -d '{"VolumeId":"'"$VOL"'","EndpointId":"'"$EP"'","LUN":5}'
```
- Poll task, then list:
```bash
curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Mappings | jq
```

3) Unmap (async)
```bash
MAP=$(curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Mappings | jq -r '.Members[0].Id')

curl -i -X DELETE http://localhost:8080/redfish/v1/StorageServices/mock/Mappings/$MAP \
  -H 'Idempotency-Key: 00000000-0000-4000-8000-000000000003'
```
- Poll task, then list again to confirm removal.

## Idempotency Notes
- `Idempotency-Key` allows safe retries for Volume create, Mapping, and Unmapping.
- Re-sending the same request with the same key returns the original job ID.

## Tuning the Simulated Driver
- Defaults are fast and reliable. To add transient failures or delays, edit `cmd/swordfish-mock/main.go`:
```go
sim.Config{MinDelay: 10 * time.Millisecond, MaxDelay: 50 * time.Millisecond, FailProb: 0}
```
- Increase `FailProb` (e.g., `0.1`) to simulate occasional `BUSY`/retry-like scenarios.

## Troubleshooting
- State is in-memory; restart clears it.
- If `go run` complains about `go.mod` `go 1.24.4`, lower it locally to your installed version (e.g., `1.22`).
- If ports are busy, change the address in `cmd/swordfish-mock/main.go`.

