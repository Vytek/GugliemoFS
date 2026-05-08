package indexer

import (
    "sync"
    "github.com/blevesearch/bleve/v2"
)

type Indexer struct {
    index bleve.Index
    mu sync.Mutex
}

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

func (i *Indexer) Index(id, content string) error {
    i.mu.Lock(); defer i.mu.Unlock()
    return i.index.Index(id, map[string]string{"content": content})
}

func (i *Indexer) Delete(id string) error {
    i.mu.Lock(); defer i.mu.Unlock()
    return i.index.Delete(id)
}
