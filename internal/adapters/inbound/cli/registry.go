package cli

import "github.com/plu9in/pdfgo/internal/domain/ports"

type Registry struct {
	handlers map[string]ports.InboundHandler
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]ports.InboundHandler)}
}

func (r *Registry) Register(h ports.InboundHandler) {
	r.handlers[h.Kind()] = h
}

func (r *Registry) Get(kind string) (ports.InboundHandler, bool) {
	h, ok := r.handlers[kind]
	return h, ok
}
