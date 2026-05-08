package extractors

import (
    "bytes"
    "github.com/dslipak/pdf"
)

type PDFExtractor struct{}

func (p *PDFExtractor) Supports(ext string) bool { return ext == ".pdf" }

func (p *PDFExtractor) Extract(path string) (string, error) {
    r, err := pdf.Open(path)
    if err != nil { return "", err }

    var buf bytes.Buffer
    for i := 1; i <= r.NumPage(); i++ {
        page := r.Page(i)
        if page.V.IsNull() { continue }
        txt, _ := page.GetPlainText(nil)
        buf.WriteString(txt)
    }
    return buf.String(), nil
}
