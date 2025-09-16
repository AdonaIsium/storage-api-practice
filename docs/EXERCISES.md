# Hands-on Exercises: Storage Orchestration, SAN, and Data Management

These 10 exercises use this mock storage API to explore common storage-appliance/SAN workflows and control‑plane patterns. They mix operational tasks (using the API) and small extensions you can implement to deepen understanding.

Prereqs: server running at http://localhost:8080 (see docs/HOWTO.md). `curl` and `jq` are helpful.

---

## 1) Provisioning 101: Create and Inspect Volumes
- Create three volumes of different sizes; note the 202 + Task pattern.
- Poll tasks to completion; list volumes and observe attributes (size, thin provisioning).
- Questions:
  - Why do many arrays return a task instead of synchronous success?
  - What metadata might you add (QoS policy, tags) and why?

Commands (example):
```bash
for i in 1 2 3; do
  curl -s -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Volumes \
    -H 'Content-Type: application/json' \
    -H "Idempotency-Key: 00000000-0000-4000-8000-00000000000$i" \
    -d "{\"Name\":\"vol$i\",\"CapacityBytes\":$((i*1073741824)),\"ThinProvisioned\":true}" -i | sed -n 's/Location: //p'
done
curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Volumes | jq
```

---

## 2) Endpoints (Hosts) and Initiators
- Create hosts with different initiator types (IQN, WWPN, NQN). Mix formats.
- Reflect on how initiator identity is used by arrays to authorize access.
- Try creating two hosts with the same identifier; observe conflict behavior vs. current mock’s behavior and consider improvements.

Example:
```bash
curl -s -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Endpoints \
  -H 'Content-Type: application/json' \
  -d '{"Name":"host-iqn","Identifiers":[{"DurableNameFormat":"iqn","DurableName":"iqn.1993-08.com.example:01:abc"}]}' -i
curl -s -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Endpoints \
  -H 'Content-Type: application/json' \
  -d '{"Name":"host-nqn","Identifiers":[{"DurableNameFormat":"nqn","DurableName":"nqn.2014-08.org.nvmexpress:uuid:1234"}]}' -i
curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Endpoints | jq
```

---

## 3) LUN Assignment and Volume Mapping
- Map two different volumes to the same host with explicit LUNs.
- Map one volume without specifying LUN and observe auto‑assignment.
- Discuss LUN collisions and host-side device discovery.

Example:
```bash
VOL=$(curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Volumes | jq -r '.Members[0].Id')
EP=$(curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Endpoints | jq -r '.Members[0].Id')
# explicit LUN
curl -s -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Mappings \
  -H 'Content-Type: application/json' -H 'Idempotency-Key: 00000000-0000-4000-8000-000000000101' \
  -d '{"VolumeId":"'"$VOL"'","EndpointId":"'"$EP"'","LUN":10}' -i
# auto LUN
curl -s -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Mappings \
  -H 'Content-Type: application/json' -H 'Idempotency-Key: 00000000-0000-4000-8000-000000000102' \
  -d '{"VolumeId":"'"$VOL"'","EndpointId":"'"$EP"'"}' -i
curl -s http://localhost:8080/redfish/v1/StorageServices/mock/Mappings | jq
```

---

## 4) Idempotency Under Retries
- Send the same create‑volume or map request multiple times with the same `Idempotency-Key`.
- Observe that the Task ID is stable and no duplicates are created.
- Change the key and compare behavior.

Example:
```bash
KEY=00000000-0000-4000-8000-00000000AA01
curl -s -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Volumes \
  -H 'Content-Type: application/json' -H "Idempotency-Key: $KEY" \
  -d '{"Name":"vol-idem","CapacityBytes":1073741824,"ThinProvisioned":true}' -i
# repeat with same key
curl -s -X POST http://localhost:8080/redfish/v1/StorageServices/mock/Volumes \
  -H 'Content-Type: application/json' -H "Idempotency-Key: $KEY" \
  -d '{"Name":"vol-idem","CapacityBytes":1073741824,"ThinProvisioned":true}' -i
```

---

## 5) Transient Failures and Backoff (Operational Resilience)
- Set `FailProb` in `cmd/swordfish-mock/main.go` (e.g., 0.2) to simulate array BUSY.
- Observe which operations fail intermittently and how tasks record errors.
- Implement a simple client retry with idempotency keys.

Questions:
- Where would exponential backoff live in a real controller? How would you detect retryable conditions?

---

## 6) Audit Trails
- Add an API endpoint to expose recent audit events from `AuditRepo`.
- Create/mapping operations, then fetch the audit log and correlate `CorrelationID`, `job_id`, and resources involved.
- Discuss the importance of audit immutability and retention policies.

Hint: follow existing Swordfish handlers in `api/swordfish/server.go`.

---

## 7) Optimistic Concurrency and CAS Updates
- Extend `VolumeRepo.UpdateCAS` usage: simulate concurrent updates to volume metadata.
- Add a PATCH endpoint for volume tags; implement CAS to avoid lost updates.
- Write a small test to prove conflicts are detected.

---

## 8) Capacity Planning and Thin Provisioning
- Modify the sim driver to track a fake pool capacity (e.g., 10 TB) and enforce `INSUFFICIENT_SPACE` when oversubscribed beyond a threshold.
- Expose capacity metrics via a `/StoragePools`-like endpoint.
- Discuss thin vs. thick provisioning and host-visible capacity.

---

## 9) Multipath & Host Grouping (Design Exercise)
- Design (and optionally prototype) host groups and volume access groups.
- Propose JSON shapes and routes for:
  - Creating groups, adding endpoints/volumes, mapping groups.
- Consider how LUN conflicts and ACLs should be handled.

---

## 10) Snapshots/Clones (Extension)
- Extend the driver/service to support snapshots and clones of volumes.
- Add async jobs and endpoints (e.g., `POST /Volumes/{id}/Snapshots`), following the same task pattern.
- Discuss metadata you would capture (parent, creation time, dependency tree) and how clones affect capacity.

---

By the end of these exercises, you will have practiced:
- Async control-plane patterns with idempotent operations
- Host/initiator identity and volume access via LUN mapping
- Failure modeling and resilience
- Auditing, concurrency control, and capacity considerations
- Designing API shapes consistent with Redfish/Swordfish conventions

