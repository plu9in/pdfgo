package document_test

import (
	"testing"

	"github.com/plu9in/pdfgo/internal/domain/document"
)

func TestGetPaperSize_KnownSizes(t *testing.T) {
	ps, err := document.GetPaperSize("A4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps.Width != 595.28 || ps.Height != 841.89 {
		t.Fatalf("unexpected A4 size: %+v", ps)
	}

	ps2, err := document.GetPaperSize("letter") // insensible à la casse
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ps2.Width != 612.0 || ps2.Height != 792.0 {
		t.Fatalf("unexpected Letter size: %+v", ps2)
	}
}

func TestGetPaperSize_Unknown(t *testing.T) {
	if _, err := document.GetPaperSize("A7"); err == nil {
		t.Fatalf("expected error for unknown size")
	}
}

func TestGetPaperSize_Empty(t *testing.T) {
	if _, err := document.GetPaperSize(""); err == nil {
		t.Fatalf("expected error for empty name")
	}
}
