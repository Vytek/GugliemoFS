package extractors

type Registry struct { list []Extractor }

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) Register(e Extractor) {
    r.list = append(r.list, e)
}

func (r *Registry) Get(ext string) Extractor {
    for _, e := range r.list {
        if e.Supports(ext) { return e }
    }
    return nil
}
