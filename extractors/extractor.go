package extractors

type Extractor interface {
    Supports(ext string) bool
    Extract(path string) (string, error)
}
