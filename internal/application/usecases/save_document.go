package usecases

import (
	"context"

	"github.com/plu9in/pdfgo/internal/domain/document"
	"github.com/plu9in/pdfgo/internal/domain/ports"
)

// SaveDocument orchestre la validation et délègue au port DocumentSaver.
type SaveDocument struct {
	saver ports.DocumentSaver
}

func NewSaveDocument(saver ports.DocumentSaver) *SaveDocument {
	return &SaveDocument{saver: saver}
}

func (uc *SaveDocument) Execute(ctx context.Context, cfg document.DocumentConfig) (string, error) {
	if err := cfg.Validate(); err != nil {
		return "", err
	}
	return uc.saver.Save(ctx, cfg)
}
