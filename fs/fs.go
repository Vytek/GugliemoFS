package fs

import (
	"context"
	"os"
	"path/filepath"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"guglielmofs/pipeline"
)

// FS implements the FUSE filesystem interface, allowing us to serve a directory structure and files.
// It uses the Pipeline to manage file processing and indexing.
type FS struct {
	RootDir  string
	Pipeline *pipeline.Pipeline
}

// Root returns the root directory node of the filesystem, which is the entry point for all file operations.
func (f *FS) Root() (fs.Node, error) {
	return &Dir{f.RootDir, f.Pipeline}, nil
}

// Dir represents a directory in the filesystem. It implements the fs.Node and fs.HandleReadDirAller interfaces.
type Dir struct {
	path string
	pipe *pipeline.Pipeline
}

// Attr retrieves the attributes of the directory, such as permissions, and fills the fuse.Attr structure.
func (d *Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Mode = os.ModeDir | 0755
	return nil
}

// ReadDirAll lists the contents of the directory, returning a slice of fuse.Dirent for each entry.
func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list, _ := os.ReadDir(d.path)
	var out []fuse.Dirent
	for _, f := range list {
		t := fuse.DT_File
		if f.IsDir() {
			t = fuse.DT_Dir
		}
		out = append(out, fuse.Dirent{Name: f.Name(), Type: t})
	}
	return out, nil
}

// Lookup finds a child node (file or directory) by name and returns the corresponding fs.Node.
func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	full := filepath.Join(d.path, name)
	fi, err := os.Stat(full)
	if err != nil {
		return nil, fuse.ENOENT
	}
	if fi.IsDir() {
		return &Dir{full, d.pipe}, nil
	}
	return &File{full, d.pipe}, nil
}

// File represents a file in the filesystem. It implements the fs.Node and fs.HandleReadAller interfaces.
type File struct {
	path string
	pipe *pipeline.Pipeline
}

// Attr retrieves the attributes of the file, such as size and permissions, and fills the fuse.Attr structure.
func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	fi, _ := os.Stat(f.path)
	a.Size = uint64(fi.Size())
	a.Mode = 0644
	return nil
}

// ReadAll reads the entire content of the file and returns it as a byte slice.
func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	return os.ReadFile(f.path)
}
