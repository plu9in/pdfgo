package pdf_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/plu9in/pdfgo/internal/adapters/outbound/pdf"
	"github.com/plu9in/pdfgo/internal/domain/document"
	"github.com/plu9in/pdfgo/internal/domain/ports"
)

func TestPDFSaver_ImplementsDocumentSaver_AndWritesA4PortraitPDF(t *testing.T) {
	tmp := t.TempDir()
	out := filepath.Join(tmp, "doc.pdf")
	cfg := document.DocumentConfig{
		Paper:       "A4",
		Orientation: document.Portrait,
		Output:      document.Output{Format: "pdf", Path: out},
	}

	// compile-time check
	var _ ports.DocumentSaver = (*pdf.PDFSaver)(nil)

	s := pdf.NewPDFSaver()
	path, err := s.Save(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Save error: %v", err)
	}
	if path != out {
		t.Fatalf("unexpected returned path: %s", path)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !strings.HasPrefix(string(data), "%PDF-1.4") {
		t.Fatalf("not a PDF-1.4 file")
	}
	if !strings.Contains(string(data), "/MediaBox [0 0 595.28 841.89]") {
		t.Fatalf("expected A4 portrait MediaBox, got:\n%s", string(data))
	}
}
