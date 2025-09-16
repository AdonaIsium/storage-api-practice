package mem

import (
    "context"
    "sync"

    "github.com/AdonaIsium/storage-api-practice/core"
    "github.com/AdonaIsium/storage-api-practice/store"
)

type MappingRepo struct {
    mu     sync.RWMutex
    byID   map[core.MappingID]*core.Mapping
    byVH   map[string]core.MappingID
    byVol  map[core.VolumeID][]core.MappingID
    byHost map[core.HostID][]core.MappingID
}

func NewMappingRepo() *MappingRepo {
    return &MappingRepo{
        byID:   make(map[core.MappingID]*core.Mapping),
        byVH:   make(map[string]core.MappingID),
        byVol:  make(map[core.VolumeID][]core.MappingID),
        byHost: make(map[core.HostID][]core.MappingID),
    }
}

func (r *MappingRepo) Save(ctx context.Context, m *core.Mapping) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    key := string(m.VolumeID) + "|" + string(m.HostID)
    if existing, ok := r.byVH[key]; ok && existing != m.ID {
        return store.ErrConflict
    }
    if prev, ok := r.byID[m.ID]; ok {
        // remove from old indexes if vol/host changed
        if prev.VolumeID != m.VolumeID || prev.HostID != m.HostID {
            // remove prev from lists
            r.removeFromLists(prev)
        }
    }
    cp := *m
    r.byID[m.ID] = &cp
    r.byVH[key] = m.ID
    r.byVol[m.VolumeID] = appendIfMissing(r.byVol[m.VolumeID], m.ID)
    r.byHost[m.HostID] = appendIfMissing(r.byHost[m.HostID], m.ID)
    return nil
}

func (r *MappingRepo) Get(ctx context.Context, id core.MappingID) (*core.Mapping, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    m, ok := r.byID[id]
    if !ok {
        return nil, store.ErrNotFound
    }
    cp := *m
    return &cp, nil
}

func (r *MappingRepo) GetByName(ctx context.Context, _ string) (*core.Mapping, error) {
    // no natural name for mapping; not found
    return nil, store.ErrNotFound
}

func (r *MappingRepo) List(ctx context.Context, f core.MappingFilter) ([]*core.Mapping, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    var ids []core.MappingID
    if f.VolumeID != "" {
        ids = append(ids, r.byVol[core.VolumeID(f.VolumeID)]...)
    }
    if f.HostID != "" {
        ids = append(ids, r.byHost[core.HostID(f.HostID)]...)
    }
    if f.VolumeID == "" && f.HostID == "" {
        // list all
        for id := range r.byID {
            ids = append(ids, id)
        }
    }
    out := make([]*core.Mapping, 0, len(ids))
    seen := map[core.MappingID]struct{}{}
    for _, id := range ids {
        if _, ok := seen[id]; ok {
            continue
        }
        seen[id] = struct{}{}
        if m, ok := r.byID[id]; ok {
            cp := *m
            out = append(out, &cp)
        }
    }
    return out, nil
}

func (r *MappingRepo) Delete(ctx context.Context, id core.MappingID) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    m, ok := r.byID[id]
    if !ok {
        return store.ErrNotFound
    }
    delete(r.byID, id)
    delete(r.byVH, string(m.VolumeID)+"|"+string(m.HostID))
    r.removeFromLists(m)
    return nil
}

func (r *MappingRepo) UpdateCAS(ctx context.Context, _ *core.Mapping, _ int64) error {
    // not used in this mock; could implement version later
    return nil
}

func appendIfMissing(s []core.MappingID, id core.MappingID) []core.MappingID {
    for _, v := range s {
        if v == id {
            return s
        }
    }
    return append(s, id)
}

func (r *MappingRepo) removeFromLists(m *core.Mapping) {
    // remove id from byVol and byHost slices
    if list, ok := r.byVol[m.VolumeID]; ok {
        r.byVol[m.VolumeID] = removeID(list, m.ID)
    }
    if list, ok := r.byHost[m.HostID]; ok {
        r.byHost[m.HostID] = removeID(list, m.ID)
    }
}

func removeID(s []core.MappingID, id core.MappingID) []core.MappingID {
    out := s[:0]
    for _, v := range s {
        if v != id {
            out = append(out, v)
        }
    }
    return out
}
