package cli

import (
	"context"

	"github.com/plu9in/pdfgo/internal/adapters/inbound/handlers"
	"github.com/plu9in/pdfgo/internal/domain/document"
)

// SaveDocumentExecutor est le contrat attendu côté CLI pour le wiring par défaut.
// Il correspond au use case réel: Execute(ctx context.Context, cfg document.DocumentConfig).
type SaveDocumentExecutor interface {
	Execute(ctx context.Context, cfg document.DocumentConfig) (string, error)
}

// saveDocUCAdapter adapte SaveDocumentExecutor -> handlers.SaveDocumentUseCase
type saveDocUCAdapter struct{ inner SaveDocumentExecutor }

// Le handler attend Execute(ctx any, cfg ...); on ignore ctx et passons Background().
func (a *saveDocUCAdapter) Execute(_ any, cfg document.DocumentConfig) (string, error) {
	return a.inner.Execute(context.Background(), cfg)
}

// NewDefaultRegistry enregistre les handlers “par défaut”, dont save_document.
func NewDefaultRegistry(saveUC SaveDocumentExecutor) *Registry {
	reg := NewRegistry()
	reg.Register(handlers.NewSaveDocumentHandler(&saveDocUCAdapter{inner: saveUC}))
	return reg
}
