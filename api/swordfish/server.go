package swordfish

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/AdonaIsium/storage-api-practice/core"
	"github.com/AdonaIsium/storage-api-practice/jobs"
	"github.com/AdonaIsium/storage-api-practice/pkg/ids"
)

type Server struct {
	deps jobs.Deps
	prov core.ProvisionService
	conn core.ConnectivityService
	jobs core.JobRepo
}

func New(deps jobs.Deps, prov core.ProvisionService, conn core.ConnectivityService) *Server {
	return &Server{deps: deps, prov: prov, conn: conn, jobs: deps.Jobs}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/redfish/v1", s.handleRoot)
	mux.HandleFunc("/redfish/v1/StorageServices", s.handleStorageServices)
	mux.HandleFunc("/redfish/v1/StorageServices/mock", s.handleStorageService)
	mux.HandleFunc("/redfish/v1/StorageServices/mock/Volumes", s.handleVolumes)
	mux.HandleFunc("/redfish/v1/StorageServices/mock/Endpoints", s.handleEndpoints)
	mux.HandleFunc("/redfish/v1/StorageServices/mock/Mappings", s.handleMappings)
	mux.HandleFunc("/redfish/v1/StorageServices/mock/Mappings/", s.handleMappingByID)
	mux.HandleFunc("/redfish/v1/TaskService/Tasks/", s.handleTask)
	return mux
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"@odata.id":       "/redfish/v1",
		"StorageServices": map[string]string{"@odata.id": "/redfish/v1/StorageServices"},
		"TaskService":     map[string]string{"@odata.id": "/redfish/v1/TaskService"},
	})
}

func (s *Server) handleStorageServices(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"Members@odata.count": 1,
		"Members":             []any{map[string]string{"@odata.id": "/redfish/v1/StorageServices/mock"}},
	})
}

func (s *Server) handleStorageService(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"Id":      "mock",
		"Name":    "Mock Storage Service",
		"Volumes": map[string]string{"@odata.id": "/redfish/v1/StorageServices/mock/Volumes"},
	})
}

type createVolumeBody struct {
	Name          string `json:"Name"`
	CapacityBytes int64  `json:"CapacityBytes"`
	Thin          bool   `json:"ThinProvisioned"`
}

func (s *Server) handleVolumes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		vols, _ := s.deps.Volumes.List(r.Context(), core.VolumeFilter{})
		members := make([]any, 0, len(vols))
		for _, v := range vols {
			members = append(members, map[string]any{
				"@odata.id":       "/redfish/v1/StorageServices/mock/Volumes/" + string(v.ID),
				"Id":              string(v.ID),
				"Name":            v.Name,
				"CapacityBytes":   v.SizeBytes,
				"ThinProvisioned": v.Thin,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"Members@odata.count": len(members), "Members": members})
	case http.MethodPost:
		var body createVolumeBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		idem := r.Header.Get("Idempotency-Key")
		if !ids.IsValid(idem) {
			idem = ids.NewIdemKey()
		}
		jobID, err := s.prov.RequestCreateVolume(r.Context(), core.CreateVolumeRequest{
			Name: body.Name, SizeBytes: body.CapacityBytes, Thin: body.Thin,
		}, idem)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", "/redfish/v1/TaskService/Tasks/"+string(jobID))
		writeJSON(w, http.StatusAccepted, map[string]any{"Task": map[string]string{"@odata.id": "/redfish/v1/TaskService/Tasks/" + string(jobID)}})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleTask(w http.ResponseWriter, r *http.Request) {
	// url: /redfish/v1/TaskService/Tasks/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/redfish/v1/TaskService/Tasks/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	id := core.JobID(parts[0])
	j, err := s.jobs.Get(r.Context(), id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	st := map[core.JobState]string{
		core.JobPending:   "Running",
		core.JobRunning:   "Running",
		core.JobSucceeded: "Completed",
		core.JobFailed:    "Exception",
	}[j.State]
	writeJSON(w, http.StatusOK, map[string]any{
		"Id":        string(j.ID),
		"TaskState": st,
		"StartTime": j.CreatedAt.Format(time.RFC3339Nano),
		"EndTime":   j.UpdatedAt.Format(time.RFC3339Nano),
		"Name":      j.Name,
		"Message":   j.ErrorMsg,
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// ---- Endpoints (Hosts) ----
type createEndpointBody struct {
	Name        string `json:"Name"`
	Identifiers []struct {
		DurableNameFormat string `json:"DurableNameFormat"`
		DurableName       string `json:"DurableName"`
	} `json:"Identifiers"`
}

func (s *Server) handleEndpoints(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		hs, _ := s.deps.Hosts.List(r.Context(), core.HostFilter{})
		members := make([]any, 0, len(hs))
		for _, h := range hs {
			members = append(members, map[string]any{"@odata.id": "/redfish/v1/StorageServices/mock/Endpoints/" + string(h.ID), "Id": string(h.ID), "Name": h.Name})
		}
		writeJSON(w, http.StatusOK, map[string]any{"Members@odata.count": len(members), "Members": members})
	case http.MethodPost:
		var body createEndpointBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		idsList := []core.HostIdentity{}
		for _, id := range body.Identifiers {
			t := core.IdentityType(strings.ToLower(id.DurableNameFormat))
			idsList = append(idsList, core.HostIdentity{Type: t, Value: id.DurableName})
		}
		jobID, err := s.conn.RequestCreateHost(r.Context(), core.CreateHostRequest{Name: body.Name, Identities: idsList})
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		w.Header().Set("Location", "/redfish/v1/TaskService/Tasks/"+string(jobID))
		writeJSON(w, http.StatusAccepted, map[string]any{"Task": map[string]string{"@odata.id": "/redfish/v1/TaskService/Tasks/" + string(jobID)}})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ---- Mappings ----
type createMappingBody struct {
	VolumeID   string `json:"VolumeId"`
	EndpointID string `json:"EndpointId"`
	LUN        *int   `json:"LUN,omitempty"`
}

func (s *Server) handleMappings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ms, _ := s.deps.Mappings.List(r.Context(), core.MappingFilter{})
		members := make([]any, 0, len(ms))
		for _, m := range ms {
			members = append(members, map[string]any{"@odata.id": "/redfish/v1/StorageServices/mock/Mappings/" + string(m.ID), "Id": string(m.ID), "VolumeId": string(m.VolumeID), "EndpointId": string(m.HostID), "LUN": m.LUN})
		}
		writeJSON(w, http.StatusOK, map[string]any{"Members@odata.count": len(members), "Members": members})
	case http.MethodPost:
		var body createMappingBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		idem := r.Header.Get("Idempotency-Key")
		if !ids.IsValid(idem) { idem = ids.NewIdemKey() }
		jobID, err := s.conn.RequestMapVolume(r.Context(), core.VolumeID(body.VolumeID), core.HostID(body.EndpointID), core.MapOptions{LUN: body.LUN}, idem)
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		w.Header().Set("Location", "/redfish/v1/TaskService/Tasks/"+string(jobID))
		writeJSON(w, http.StatusAccepted, map[string]any{"Task": map[string]string{"@odata.id": "/redfish/v1/TaskService/Tasks/" + string(jobID)}})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleMappingByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/redfish/v1/StorageServices/mock/Mappings/")
	if id == "" { http.NotFound(w, r); return }
	switch r.Method {
	case http.MethodDelete:
		idem := r.Header.Get("Idempotency-Key")
		if !ids.IsValid(idem) { idem = ids.NewIdemKey() }
		jobID, err := s.conn.RequestUnmap(r.Context(), core.MappingID(id), idem)
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		w.Header().Set("Location", "/redfish/v1/TaskService/Tasks/"+string(jobID))
		writeJSON(w, http.StatusAccepted, map[string]any{"Task": map[string]string{"@odata.id": "/redfish/v1/TaskService/Tasks/" + string(jobID)}})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
