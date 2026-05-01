# GugliemoFS

![GugliemoFS](GuglielmoFS.png)

**GugliemoFS** is a distributed, FUSE-based filesystem designed to index and expose every type of document as a navigable filesystem tree. Built on top of [XtreemFS](http://www.xtreemfs.org/) and [Apache Solr](http://lucene.apache.org/solr/), it powers the **Tera Document Management System** by making full-text search results browsable as ordinary directories and files.

---

## Overview

Traditional document management systems store files in flat or weakly structured repositories. GugliemoFS takes a different approach: it mounts a virtual filesystem where folders represent search categories, metadata facets, or query results, and files are the actual documents — all accessible through any standard POSIX application without modification.

Under the hood, every file operation (open, read, listdir) triggers a live query against the Solr search index, so the filesystem view is always up to date with the document corpus.

---

## Features

- **FUSE-based virtual filesystem** — mount anywhere on Linux/macOS with no kernel patches required
- **Universal document indexing** — supports PDF, Office documents, plain text, images with EXIF metadata, email archives, and any other format Solr can extract content from (via Apache Tika)
- **Faceted directory tree** — browse documents by type, author, date, tag, or any indexed field as if they were nested folders
- **Full-text search as a path** — type a search query as a directory path and its results appear as files
- **Distributed storage** — files are physically stored and replicated across XtreemFS object storage nodes
- **Transparent access** — standard tools (`ls`, `cp`, `find`, `grep`, editors) work without any changes
- **Real-time index updates** — new documents are indexed and immediately visible in the filesystem

---

## Architecture

```
┌─────────────────────────────────────────────┐
│              POSIX Applications              │
│       (ls, cat, find, editors, ...)          │
└───────────────────┬─────────────────────────┘
                    │  POSIX calls
┌───────────────────▼─────────────────────────┐
│              GugliemoFS (FUSE)               │
│  ┌─────────────────┐  ┌────────────────────┐│
│  │  Virtual FS     │  │  Metadata Cache    ││
│  │  (directory/    │  │  (TTL-based,       ││
│  │   file mapping) │  │   query results)   ││
│  └────────┬────────┘  └────────────────────┘│
└───────────┼─────────────────────────────────┘
            │
  ┌─────────▼──────────┐    ┌─────────────────┐
  │    Apache Solr     │    │   XtreemFS      │
  │  (Full-text index, │    │ (Distributed    │
  │   facets, metadata)│    │  object store)  │
  └────────────────────┘    └─────────────────┘
            │                        │
  ┌─────────▼────────────────────────▼────────┐
  │            Document Corpus                 │
  │  (PDF, DOCX, TXT, images, emails, ...)    │
  └────────────────────────────────────────────┘
```

---

## Directory Layout

When mounted, GugliemoFS exposes the following virtual directory structure:

```
/mnt/gugliemofs/
├── by-type/
│   ├── pdf/
│   ├── docx/
│   ├── txt/
│   └── ...
├── by-author/
│   ├── mario.rossi/
│   └── ...
├── by-date/
│   ├── 2024/
│   │   ├── 01/
│   │   └── ...
│   └── ...
├── by-tag/
│   ├── invoice/
│   ├── contract/
│   └── ...
└── search/
    └── <query>/          ← results appear here as files
```

---

## Requirements

- Linux kernel ≥ 2.6.14 with FUSE support (or macOS with [macFUSE](https://osxfuse.github.io/))
- [libfuse](https://github.com/libfuse/libfuse) ≥ 3.x
- [Apache Solr](http://lucene.apache.org/solr/) ≥ 8.x (with Apache Tika content extraction)
- [XtreemFS](http://www.xtreemfs.org/) object storage cluster
- Java 11+ (required by Solr)

---

## Getting Started

### 1. Start XtreemFS

Follow the [XtreemFS quickstart guide](http://www.xtreemfs.org/xtfs-guide-1.5.1/index.html) to bring up a DIR, MRC, and at least one OSD node.

### 2. Start Apache Solr

```bash
bin/solr start
bin/solr create -c gugliemofs
```

Configure the `gugliemofs` core with the provided `schema.xml` to enable full-text extraction and faceting.

### 3. Build GugliemoFS

```bash
git clone https://github.com/Vytek/GugliemoFS.git
cd GugliemoFS
make
```

### 4. Mount the filesystem

```bash
mkdir /mnt/gugliemofs
./gugliemofs \
  --solr-url http://localhost:8983/solr/gugliemofs \
  --xtreemfs-url pbrpc://localhost/gugliemofs-vol \
  /mnt/gugliemofs
```

### 5. Index documents

Place documents into the XtreemFS volume (or use the provided ingestion tool). GugliemoFS will automatically submit them to Solr for indexing.

```bash
cp /path/to/documents/*.pdf /mnt/gugliemofs/by-type/pdf/
```

### 6. Browse and search

```bash
ls /mnt/gugliemofs/by-author/mario.rossi/
ls "/mnt/gugliemofs/search/annual report 2023/"
cat "/mnt/gugliemofs/search/annual report 2023/report_final.pdf"
```

---

## Configuration

| Option | Default | Description |
|---|---|---|
| `--solr-url` | `http://localhost:8983/solr/gugliemofs` | Solr core endpoint |
| `--xtreemfs-url` | — | XtreemFS volume URL |
| `--cache-ttl` | `30` | Metadata cache TTL in seconds |
| `--facets` | `type,author,date,tag` | Comma-separated list of facet fields to expose as directories |
| `--log-level` | `info` | Logging verbosity (`debug`, `info`, `warn`, `error`) |

---

## Use Cases

- **Enterprise Document Management** — mount the entire document repository as a filesystem and use existing backup, search, and processing tools directly
- **Legal & Compliance** — browse contracts and invoices by date and tag without learning a new UI
- **Research Archives** — navigate large corpora of papers by author, year, or keyword
- **Automated Pipelines** — use shell scripts or cron jobs to process new documents as they appear in the virtual filesystem

---

## Contributing

Pull requests and issue reports are welcome. Please open an issue to discuss major changes before submitting a PR.

---

## License

This project is released under the [Apache License 2.0](LICENSE).

---

## Related Projects

- [XtreemFS](http://www.xtreemfs.org/) — distributed filesystem used as the physical storage backend
- [Apache Solr](http://lucene.apache.org/solr/) — search platform used for indexing and querying documents
- [Apache Tika](https://tika.apache.org/) — content extraction toolkit integrated via Solr Cell
- [libfuse](https://github.com/libfuse/libfuse) — FUSE userspace library

