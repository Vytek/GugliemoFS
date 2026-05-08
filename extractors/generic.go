package extractors

import (
	"strings"

	documentcrack "github.com/taigrr/document-crack"
)

type GenericExtractor struct{}

func (g *GenericExtractor) Supports(ext string) bool {
	switch ext {
	case ".txt", ".doc", ".docx", ".ods":
		return true
	}
	return false
}

func (g *GenericExtractor) Extract(path string) (string, error) {
	doc, err := documentcrack.FromFile(path)
	if err != nil {
		return "", err
	}
	return strings.Join(doc.Content, "\n"), nil
}
