package extractors

import "github.com/taigrr/document-crack/crack"

type GenericExtractor struct{}

func (g *GenericExtractor) Supports(ext string) bool {
    switch ext {
    case ".txt",".doc",".docx",".ods":
        return true
    }
    return false
}

func (g *GenericExtractor) Extract(path string) (string, error) {
    return crack.ExtractText(path)
}
