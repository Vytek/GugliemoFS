package pipeline

import (
    "log"
    "path/filepath"

    "fuse-indexer/extractors"
    "fuse-indexer/indexer"
    "fuse-indexer/utils"
)

type Pipeline struct {
    reg *extractors.Registry
    idx *indexer.Indexer
    meta *indexer.MetaStore
}

func NewPipeline(r *extractors.Registry, i *indexer.Indexer, m *indexer.MetaStore) *Pipeline {
    return &Pipeline{r,i,m}
}

func (p *Pipeline) Process(path string) {
    ext := filepath.Ext(path)
    ex := p.reg.Get(ext)
    if ex == nil { return }

    hash, err := utils.FileHash(path)
    if err != nil { return }

    if hash == p.meta.Get(path) { return }

    log.Println("Index:", path)

    txt, err := ex.Extract(path)
    if err != nil { return }

    if err := p.idx.Index(path, txt); err == nil {
        p.meta.Set(path, hash)
    }
}

func (p *Pipeline) Delete(path string) {
    p.idx.Delete(path)
    p.meta.Delete(path)
}

func (p *Pipeline) Rename(old, new string) {
    p.meta.Rename(old,new)
}
