package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/blevesearch/bleve/v2"
	"github.com/dslipak/pdf"
)

var (
	backingDir string // Dove vengono salvati fisicamente i file
	mountPoint string // La cartella virtuale FUSE
	searchIndex bleve.Index
)

// Document rappresenta la struttura del dato indicizzato in Bleve
type PDFDocument struct {
	Path    string
	Content string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Uso: %s <backing_dir> <mount_point>\n", os.Args[0])
		os.Exit(1)
	}

	backingDir = os.Args[1]
	mountPoint = os.Args[2]

	// 1. Setup della directory di supporto
	if err := os.MkdirAll(backingDir, 0755); err != nil {
		log.Fatalf("Errore nel creare la backing dir: %v", err)
	}

	// 2. Setup dell'indice Bleve
	indexPath := "pdf_index.bleve"
	var err error
	if searchIndex, err = bleve.Open(indexPath); err != nil {
		log.Println("Indice non trovato, ne creo uno nuovo...")
		mapping := bleve.NewIndexMapping()
		searchIndex, err = bleve.New(indexPath, mapping)
		if err != nil {
			log.Fatalf("Errore nella creazione dell'indice Bleve: %v", err)
		}
	}
	defer searchIndex.Close()

	// 3. Mount del Filesystem FUSE
	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("pdf_indexer"),
		fuse.Subtype("pdffs"),
	)
	if err != nil {
		log.Fatalf("Errore nel mount FUSE: %v", err)
	}
	defer c.Close()
	defer fuse.Unmount(mountPoint)

	log.Printf("Filesystem montato su %s (Backing su %s)", mountPoint, backingDir)
	log.Println("Copia un file .pdf nella cartella di mount per testarlo!")

	// 4. Avvio del server FUSE
	err = fs.Serve(c, AppFS{})
	if err != nil {
		log.Fatalf("Errore durante l'esecuzione del filesystem: %v", err)
	}
}

// ==========================================
// Logica FUSE: Strutture Base
// ==========================================

// AppFS è la radice del filesystem
type AppFS struct{}

func (AppFS) Root() (fs.Node, error) {
	return Dir{Path: backingDir}, nil
}

// Dir rappresenta una directory nel nostro FUSE
type Dir struct {
	Path string
}

func (d Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	stat, err := os.Stat(d.Path)
	if err != nil {
		return err
	}
	a.Inode = 1 // Semplificazione: inode fittizio
	a.Mode = os.ModeDir | stat.Mode()
	return nil
}

// Lookup trova un file o cartella all'interno della directory
func (d Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	fullPath := filepath.Join(d.Path, name)
	stat, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fuse.ENOENT
		}
		return nil, err
	}

	if stat.IsDir() {
		return Dir{Path: fullPath}, nil
	}
	return &File{Path: fullPath}, nil
}

// Create intercetta la creazione di un nuovo file (es. quando copi un PDF)
func (d Dir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	fullPath := filepath.Join(d.Path, req.Name)
	f, err := os.OpenFile(fullPath, int(req.Flags), req.Mode)
	if err != nil {
		return nil, nil, err
	}
	
	fileNode := &File{
		Path: fullPath,
		fd:   f,
	}
	return fileNode, fileNode, nil
}

// ==========================================
// Logica FUSE: Gestione File
// ==========================================

// File rappresenta un file nel filesystem FUSE
type File struct {
	Path string
	fd   *os.File
	mu   sync.Mutex
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	stat, err := os.Stat(f.Path)
	if err != nil {
		return err
	}
	a.Mode = stat.Mode()
	a.Size = uint64(stat.Size())
	return nil
}

// Write gestisce la scrittura fisica dei dati nella backing directory
func (f *File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.fd == nil {
		file, err := os.OpenFile(f.Path, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		f.fd = file
	}

	n, err := f.fd.WriteAt(req.Data, req.Offset)
	resp.Size = n
	return err
}

// Release viene chiamato quando il file descriptor viene chiuso (es. fine della copia)
func (f *File) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	f.mu.Lock()
	if f.fd != nil {
		f.fd.Close()
		f.fd = nil
	}
	f.mu.Unlock()

	// TRIGGER: Se è un PDF, avviamo l'indicizzazione in background
	if strings.HasSuffix(strings.ToLower(f.Path), ".pdf") {
		log.Printf("Chiusura del file PDF rilevata, avvio indicizzazione: %s", f.Path)
		go indexPDF(f.Path)
	}

	return nil
}

// ==========================================
// Logica di Parsing PDF e Indicizzazione (Bleve)
// ==========================================

func indexPDF(filePath string) {
	// 1. Apriamo il PDF con dslipak/pdf
	r, err := pdf.Open(filePath)
	if err != nil {
		log.Printf("[Errore] Impossibile aprire PDF %s: %v", filePath, err)
		return
	}

	// 2. Estraiamo il testo puro
	reader, err := r.GetPlainText()
	if err != nil {
		log.Printf("[Errore] Impossibile estrarre testo da %s: %v", filePath, err)
		return
	}
	
	var buf bytes.Buffer
	_, err = io.Copy(&buf, reader)
	if err != nil {
		log.Printf("[Errore] Errore di lettura testo da %s: %v", filePath, err)
		return
	}
	extractedText := buf.String()

	// 3. Prepariamo il documento per Bleve
	doc := PDFDocument{
		Path:    filePath,
		Content: extractedText,
	}

	// Usiamo il percorso del file come ID univoco per l'indice
	docID := filepath.Base(filePath)

	// 4. Inseriamo il documento nell'indice
	err = searchIndex.Index(docID, doc)
	if err != nil {
		log.Printf("[Errore Bleve] Fallita indicizzazione di %s: %v", filePath, err)
	} else {
		log.Printf("✔ File %s indicizzato con successo in Bleve!", docID)
	}
}