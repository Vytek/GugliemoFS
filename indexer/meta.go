package indexer

import (
    "encoding/json"
    "os"
    "sync"
)

type MetaStore struct {
    path string
    data map[string]string
    mu   sync.Mutex
}

func NewMetaStore(path string) (*MetaStore, error) {
    ms := &MetaStore{path: path, data: make(map[string]string)}

    if b, err := os.ReadFile(path); err == nil {
        json.Unmarshal(b, &ms.data)
    }

    return ms, nil
}

func (m *MetaStore) Get(p string) string {
    m.mu.Lock(); defer m.mu.Unlock()
    return m.data[p]
}

func (m *MetaStore) Set(p, h string) {
    m.mu.Lock(); defer m.mu.Unlock()
    m.data[p] = h
    m.save()
}

func (m *MetaStore) Delete(p string) {
    m.mu.Lock(); defer m.mu.Unlock()
    delete(m.data, p)
    m.save()
}

func (m *MetaStore) Rename(old, new string) {
    m.mu.Lock(); defer m.mu.Unlock()
    if h, ok := m.data[old]; ok {
        m.data[new] = h
        delete(m.data, old)
    }
    m.save()
}

func (m *MetaStore) save() {
    b, _ := json.MarshalIndent(m.data, "", "  ")
    os.WriteFile(m.path, b, 0644)
}
