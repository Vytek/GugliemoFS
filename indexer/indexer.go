package indexer

import (
	"sync"

	"github.com/blevesearch/bleve/v2"
)

// Indexer manages the Bleve index and provides thread-safe methods for indexing and deleting documents.
type Indexer struct {
	index bleve.Index
	mu    sync.Mutex
}

// NewIndexer creates a new Indexer instance with the given index path. If the index does not exist, it will be created.
func NewIndexer(path string) (*Indexer, error) {
	mapping := bleve.NewIndexMapping()
	idx, err := bleve.Open(path)
	if err != nil {
		idx, err = bleve.New(path, mapping)
		if err != nil {
			return nil, err
		}
	}
	return &Indexer{index: idx}, nil
}

// Index indexes a document with the given ID and content.
func (i *Indexer) Index(id, content string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.index.Index(id, map[string]string{"content": content})
}

// Delete removes a document with the given ID from the index.
func (i *Indexer) Delete(id string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.index.Delete(id)
}
