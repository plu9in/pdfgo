package ports

import (
	"context"

	"github.com/plu9in/pdfgo/internal/domain/document"
)

// DocumentSaver est le port de sortie pour persister/rendre un document
// selon sa configuration (PDF, PNG, etc.). Retourne un chemin (ou id) résultat.
type DocumentSaver interface {
	Save(ctx context.Context, cfg document.DocumentConfig) (string, error)
}
