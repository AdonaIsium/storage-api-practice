package mem

import (
    "context"
    "sync"

    "github.com/AdonaIsium/storage-api-practice/core"
    "github.com/AdonaIsium/storage-api-practice/store"
)

type HostRepo struct {
	mu         sync.RWMutex
	byID       map[core.HostID]*core.Host
	byName     map[string]core.HostID
	byIdentity map[string]core.HostID
}

func NewHostRepo() *HostRepo {
    return &HostRepo{
        byID:       make(map[core.HostID]*core.Host),
        byName:     make(map[string]core.HostID),
        byIdentity: make(map[string]core.HostID),
    }
}

func (r *HostRepo) Save(ctx context.Context, h *core.Host) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if h == nil {
        return nil
    }
    if existing, ok := r.byName[h.Name]; ok && existing != h.ID {
        return store.ErrConflict
    }
    // enforce identity uniqueness
    for _, id := range h.Identities {
        key := string(id.Type) + ":" + id.Value
        if existing, ok := r.byIdentity[key]; ok && existing != h.ID {
            return store.ErrConflict
        }
    }
    // update indexes
    if prev, ok := r.byID[h.ID]; ok {
        if prev.Name != h.Name {
            delete(r.byName, prev.Name)
        }
        for _, id := range prev.Identities {
            delete(r.byIdentity, string(id.Type)+":"+id.Value)
        }
    }
    cp := *h
    r.byID[h.ID] = &cp
    r.byName[h.Name] = h.ID
    for _, id := range h.Identities {
        r.byIdentity[string(id.Type)+":"+id.Value] = h.ID
    }
    return nil
}

func (r *HostRepo) Get(ctx context.Context, id core.HostID) (*core.Host, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    h, ok := r.byID[id]
    if !ok {
        return nil, store.ErrNotFound
    }
    cp := *h
    return &cp, nil
}

func (r *HostRepo) GetByName(ctx context.Context, name string) (*core.Host, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    id, ok := r.byName[name]
    if !ok {
        return nil, store.ErrNotFound
    }
    h, ok := r.byID[id]
    if !ok {
        return nil, store.ErrNotFound
    }
    cp := *h
    return &cp, nil
}

func (r *HostRepo) List(ctx context.Context, f core.HostFilter) ([]*core.Host, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    out := []*core.Host{}
    if f.NameEquals != "" {
        if id, ok := r.byName[f.NameEquals]; ok {
            if h, ok2 := r.byID[id]; ok2 {
                cp := *h
                out = append(out, &cp)
            }
        }
        return out, nil
    }
    if f.Identity != "" {
        if id, ok := r.byIdentity[f.Identity]; ok {
            if h, ok2 := r.byID[id]; ok2 {
                cp := *h
                out = append(out, &cp)
            }
        }
        return out, nil
    }
    for _, h := range r.byID {
        cp := *h
        out = append(out, &cp)
    }
    return out, nil
}

func (r *HostRepo) Delete(ctx context.Context, id core.HostID) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    h, ok := r.byID[id]
    if !ok {
        return store.ErrNotFound
    }
    delete(r.byID, id)
    delete(r.byName, h.Name)
    for _, iden := range h.Identities {
        delete(r.byIdentity, string(iden.Type)+":"+iden.Value)
    }
    return nil
}

func (r *HostRepo) UpdateCAS(ctx context.Context, h *core.Host, expectedVersion int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    cur, ok := r.byID[h.ID]
    if !ok {
        return store.ErrNotFound
    }
    if cur.UpdatedAt.UnixNano() != expectedVersion {
        return store.ErrConflict
    }
    // re-use Save path to handle indexes
    return r.Save(ctx, h)
}
