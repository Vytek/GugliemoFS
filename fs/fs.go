package fs

import (
    "bazil.org/fuse"
    "bazil.org/fuse/fs"
    "context"
    "os"
    "path/filepath"

    "fuse-indexer/pipeline"
)

type FS struct {
    RootDir string
    Pipeline *pipeline.Pipeline
}

func (f *FS) Root() (fs.Node, error) {
    return &Dir{f.RootDir,f.Pipeline}, nil
}

type Dir struct {
    path string
    pipe *pipeline.Pipeline
}

func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
    a.Mode = os.ModeDir|0755
    return nil
}

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
    list,_ := os.ReadDir(d.path)
    var out []fuse.Dirent
    for _,f := range list {
        t := fuse.DT_File
        if f.IsDir(){t=fuse.DT_Dir}
        out=append(out,fuse.Dirent{Name:f.Name(),Type:t})
    }
    return out,nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
    full := filepath.Join(d.path,name)
    fi,err:=os.Stat(full)
    if err!=nil{return nil,fuse.ENOENT}
    if fi.IsDir(){return &Dir{full,d.pipe},nil}
    return &File{full,d.pipe},nil
}

type File struct {
    path string
    pipe *pipeline.Pipeline
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
    fi,_:=os.Stat(f.path)
    a.Size=uint64(fi.Size())
    a.Mode=0644
    return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
    return os.ReadFile(f.path)
}
