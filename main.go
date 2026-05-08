package main

import (
	"log"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"guglielmofs/extractors"
	myfs "guglielmofs/fs"
	"guglielmofs/indexer"
	"guglielmofs/pipeline"
	"guglielmofs/watcher"
)

func main() {
	mount := "/mnt/fusefs"

	c, err := fuse.Mount(mount)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	reg := extractors.NewRegistry()
	reg.Register(&extractors.PDFExtractor{})
	reg.Register(&extractors.GenericExtractor{})

	idx, err := indexer.NewIndexer("index.bleve")
	if err != nil {
		log.Fatal(err)
	}

	meta, err := indexer.NewMetaStore("meta.json")
	if err != nil {
		log.Fatal(err)
	}

	pipe := pipeline.NewPipeline(reg, idx, meta)
	defer pipe.Stop()

	watcher.Start("./data", watcher.Handler{
		OnWrite:  pipe.Submit,
		OnDelete: pipe.Delete,
	})

	filesystem := &myfs.FS{
		RootDir: "./data",
	}

	if err := fs.Serve(c, filesystem); err != nil {
		log.Fatal(err)
	}
}
