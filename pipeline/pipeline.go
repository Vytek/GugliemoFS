package pipeline

import (
	"log"
	"path/filepath"

	"github.com/alitto/pond"

	"guglielmofs/extractors"
	"guglielmofs/indexer"
	"guglielmofs/utils"
)

// Pipeline manages the flow of file processing: watching for changes, extracting text, and indexing.
type Pipeline struct {
	reg  *extractors.Registry
	idx  *indexer.Indexer
	meta *indexer.MetaStore
	pool *pond.WorkerPool
}

// NewPipeline creates a new Pipeline instance with the given extractor registry, indexer, and meta store.
func NewPipeline(r *extractors.Registry, i *indexer.Indexer, m *indexer.MetaStore) *Pipeline {
	return &Pipeline{
		reg:  r,
		idx:  i,
		meta: m,
		pool: pond.New(4, 100),
	}
}

// Submit a file for processing
func (p *Pipeline) Submit(path string) {
	p.pool.Submit(func() {
		p.process(path)
	})
}

// Internal processing logic
func (p *Pipeline) process(path string) {
	ext := filepath.Ext(path)

	ex := p.reg.Get(ext)
	if ex == nil {
		return
	}

	// Check if file has changed since last indexing
	hash, err := utils.FileHash(path)
	if err != nil {
		return
	}

	// If hash matches, skip re-indexing
	if hash == p.meta.Get(path) {
		return
	}

	log.Println("Index:", path)

	// Extract text and index
	txt, err := ex.Extract(path)
	if err != nil {
		return
	}

	if p.idx.Index(path, txt) == nil {
		p.meta.Set(path, hash)
	}
}

// Delete a file from the index
func (p *Pipeline) Delete(path string) {
	p.pool.Submit(func() {
		p.idx.Delete(path)
		p.meta.Delete(path)
	})
}

// Stop the pipeline and wait for all tasks to complete
func (p *Pipeline) Stop() {
	p.pool.StopAndWait()
}
