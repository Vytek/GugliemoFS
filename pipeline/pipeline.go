package pipeline

import (
	"log"
	"path/filepath"

	"github.com/alitto/pond"

	"guglielmofs/extractors"
	"guglielmofs/indexer"
	"guglielmofs/utils"
)

type Pipeline struct {
	reg  *extractors.Registry
	idx  *indexer.Indexer
	meta *indexer.MetaStore
	pool *pond.WorkerPool
}

func NewPipeline(r *extractors.Registry, i *indexer.Indexer, m *indexer.MetaStore) *Pipeline {
	return &Pipeline{
		reg:  r,
		idx:  i,
		meta: m,
		pool: pond.New(4, 100),
	}
}

func (p *Pipeline) Submit(path string) {
	p.pool.Submit(func() {
		p.process(path)
	})
}

func (p *Pipeline) process(path string) {
	ext := filepath.Ext(path)

	ex := p.reg.Get(ext)
	if ex == nil {
		return
	}

	hash, err := utils.FileHash(path)
	if err != nil {
		return
	}

	if hash == p.meta.Get(path) {
		return
	}

	log.Println("Index:", path)

	txt, err := ex.Extract(path)
	if err != nil {
		return
	}

	if p.idx.Index(path, txt) == nil {
		p.meta.Set(path, hash)
	}
}

func (p *Pipeline) Delete(path string) {
	p.pool.Submit(func() {
		p.idx.Delete(path)
		p.meta.Delete(path)
	})
}

func (p *Pipeline) Stop() {
	p.pool.StopAndWait()
}
