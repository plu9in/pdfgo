package handlers

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/plu9in/pdfgo/internal/domain/document"
	"github.com/plu9in/pdfgo/internal/domain/ports"
)

// Vérif compile-time : SaveDocumentHandler implémente bien le port d'entrée.
var _ ports.InboundHandler = (*SaveDocumentHandler)(nil)

// SaveDocumentUseCase est le contrat attendu par le handler.
// (ctx en 'any' pour coller au fake des tests ; on pourra
// le resserrer vers context.Context plus tard.)
type SaveDocumentUseCase interface {
	Execute(ctx any, cfg document.DocumentConfig) (string, error)
}

type SaveDocumentHandler struct {
	uc SaveDocumentUseCase
}

func NewSaveDocumentHandler(uc SaveDocumentUseCase) *SaveDocumentHandler {
	return &SaveDocumentHandler{uc: uc}
}

func (h *SaveDocumentHandler) Kind() string { return "save_document" }

// Requête YAML attendue :
// kind: save_document
// spec: { ...DocumentConfig... }
type saveDocumentRequest struct {
	Kind string                   `yaml:"kind"`
	Spec *document.DocumentConfig `yaml:"spec"`
}

func (h *SaveDocumentHandler) Handle(yml []byte) error {
	var req saveDocumentRequest
	if err := yaml.Unmarshal(yml, &req); err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}
	if req.Spec == nil {
		return fmt.Errorf("missing 'spec' section")
	}

	// Validation côté handler (les tests exigent que l'UC ne soit pas appelée si invalide)
	cfg := *req.Spec
	if err := cfg.Validate(); err != nil {
		return err
	}

	// Appel du use case
	_, err := h.uc.Execute(context.Background(), cfg)
	return err
}
