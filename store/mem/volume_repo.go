package mem

import (
    "context"
    "sync"

    "github.com/AdonaIsium/storage-api-practice/core"
    "github.com/AdonaIsium/storage-api-practice/store"
)

type VolumeRepo struct {
    mu     sync.RWMutex
    byID   map[core.VolumeID]*core.Volume
    byName map[string]core.VolumeID
}

func NewVolumeRepo() *VolumeRepo {
    return &VolumeRepo{byID: make(map[core.VolumeID]*core.Volume), byName: make(map[string]core.VolumeID)}
}

func (r *VolumeRepo) Save(ctx context.Context, v *core.Volume) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if v == nil {
        return nil
    }
    if _, ok := r.byName[v.Name]; ok {
        // if same ID, allow update; if different, conflict
        if r.byName[v.Name] != v.ID {
            return store.ErrConflict
        }
    }
    // copy to avoid external mutation
    cp := *v
    r.byID[v.ID] = &cp
    r.byName[v.Name] = v.ID
    return nil
}

func (r *VolumeRepo) GetByName(ctx context.Context, name string) (*core.Volume, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    id, ok := r.byName[name]
    if !ok {
        return nil, store.ErrNotFound
    }
    v, ok := r.byID[id]
    if !ok {
        return nil, store.ErrNotFound
    }
    cp := *v
    return &cp, nil
}

func (r *VolumeRepo) List(ctx context.Context, f core.VolumeFilter) ([]*core.Volume, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    out := []*core.Volume{}
    if f.NameEquals != "" {
        if id, ok := r.byName[f.NameEquals]; ok {
            if v, ok2 := r.byID[id]; ok2 {
                cp := *v
                out = append(out, &cp)
            }
        }
        return out, nil
    }
    for _, v := range r.byID {
        if f.TagKey != "" {
            if v.Tags == nil {
                continue
            }
            val, ok := v.Tags[f.TagKey]
            if !ok {
                continue
            }
            if f.TagValue != "" && f.TagValue != val {
                continue
            }
        }
        cp := *v
        out = append(out, &cp)
    }
    return out, nil
}

func (r *VolumeRepo) Delete(ctx context.Context, id core.VolumeID) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    v, ok := r.byID[id]
    if !ok {
        return store.ErrNotFound
    }
    delete(r.byID, id)
    delete(r.byName, v.Name)
    return nil
}

func (r *VolumeRepo) UpdateCAS(ctx context.Context, v *core.Volume, expectedVersion int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    cur, ok := r.byID[v.ID]
    if !ok {
        return store.ErrNotFound
    }
    if cur.Version != expectedVersion {
        return store.ErrConflict
    }
    // update name index if changed
    if cur.Name != v.Name {
        if _, exists := r.byName[v.Name]; exists && r.byName[v.Name] != v.ID {
            return store.ErrConflict
        }
        delete(r.byName, cur.Name)
        r.byName[v.Name] = v.ID
    }
    cp := *v
    r.byID[v.ID] = &cp
    return nil
}

func (r *VolumeRepo) Get(ctx context.Context, id core.VolumeID) (*core.Volume, error) { // helper for interface completeness
    r.mu.RLock()
    defer r.mu.RUnlock()
    v, ok := r.byID[id]
    if !ok {
        return nil, store.ErrNotFound
    }
    cp := *v
    return &cp, nil
}
