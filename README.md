# GugliemoFS

![GugliemoFS](GuglielmoFS.png)

**GugliemoFS** is a distributed, FUSE-based filesystem designed to index and expose every type of document as a navigable filesystem tree. 

---

## Overview

Traditional document management systems store files in flat or weakly structured repositories. GugliemoFS takes a different approach: it mounts a virtual filesystem where folders represent search categories, metadata facets, or query results, and files are the actual documents — all accessible through any standard POSIX application without modification.


---

## Features

- **FUSE-based virtual filesystem** — mount anywhere on Linux/macOS with no kernel patches required
- **Universal document indexing** — supports PDF, Office documents, plain text, images with EXIF metadata, email archives, 
- **Faceted directory tree** — browse documents by type, author, date, tag, or any indexed field as if they were nested folders
- **Full-text search as a path** — type a search query as a directory path and its results appear as files
- **Transparent access** — standard tools (`ls`, `cp`, `find`, `grep`, editors) work without any changes
- **Real-time index updates** — new documents are indexed and immediately visible in the filesystem

---

## Architecture


---

## Requirements

- Linux kernel ≥ 2.6.14 with FUSE support (or macOS with [macFUSE](https://osxfuse.github.io/))

---

## Getting Started

### 1. Build GugliemoFS

```bash
git clone https://github.com/Vytek/GugliemoFS.git
cd GugliemoFS
make
```

### 2. Mount the filesystem

```bash
mkdir /mnt/gugliemofs
./gugliemofs \
  /mnt/gugliemofs
```

### 3. Index documents

Place documents into ... volume (or use the provided ingestotion ol). GugliemoFS will automatically submit them to Solr for indexing.

```bash
cp /path/to/documents/*.pdf /mnt/gugliemofs/by-type/pdf/
```

### 4. Browse and search

```bash
ls /mnt/gugliemofs/by-author/mario.rossi/
ls "/mnt/gugliemofs/search/annual report 2023/"
cat "/mnt/gugliemofs/search/annual report 2023/report_final.pdf"
```

---

## Configuration

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


