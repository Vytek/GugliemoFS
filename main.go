package main

import (
    "log"

    "bazil.org/fuse"
    "bazil.org/fuse/fs"

    "fuse-indexer/extractors"
    "fuse-indexer/indexer"
    "fuse-indexer/pipeline"
    "fuse-indexer/watcher"
    myfs "fuse-indexer/fs"
)

func main() {
    c, err := fuse.Mount("/mnt/fusefs")
    if err != nil { log.Fatal(err) }

    reg := extractors.NewRegistry()
    reg.Register(&extractors.PDFExtractor{})
    reg.Register(&extractors.GenericExtractor{})

    idx, _ := indexer.NewIndexer("index.bleve")
    meta, _ := indexer.NewMetaStore("meta.json")

    pipe := pipeline.NewPipeline(reg, idx, meta)

    watcher.Start("./data", watcher.Handler{
        OnWrite: pipe.Process,
        OnDelete: pipe.Delete,
    })

    fs.Serve(c, &myfs.FS{RootDir:"./data", Pipeline:pipe})
}
