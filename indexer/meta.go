package indexer

import (
	"encoding/json"
	"os"
	"sync"
)

// MetaStore manages file metadata, such as hashes, to optimize indexing by avoiding re-indexing unchanged files.
type MetaStore struct {
	path string
	data map[string]string
	mu   sync.Mutex
}

// NewMetaStore creates a new MetaStore instance with the given file path for storing metadata. It loads existing metadata if the file exists.
func NewMetaStore(path string) (*MetaStore, error) {
	ms := &MetaStore{path: path, data: make(map[string]string)}

	if b, err := os.ReadFile(path); err == nil {
		json.Unmarshal(b, &ms.data)
	}

	return ms, nil
}

// Get retrieves the hash for a given file path.
func (m *MetaStore) Get(p string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[p]
}

// Set updates the hash for a given file path and saves the metadata.
func (m *MetaStore) Set(p, h string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[p] = h
	m.save()
}

// Delete removes the metadata for a given file path and saves the metadata.
func (m *MetaStore) Delete(p string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, p)
	m.save()
}

// Rename updates the metadata for a file that has been renamed and saves the metadata.
func (m *MetaStore) Rename(old, new string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if h, ok := m.data[old]; ok {
		m.data[new] = h
		delete(m.data, old)
	}
	m.save()
}

// save writes the current metadata to the file in JSON format.
func (m *MetaStore) save() {
	b, _ := json.MarshalIndent(m.data, "", "  ")
	os.WriteFile(m.path, b, 0644)
}
