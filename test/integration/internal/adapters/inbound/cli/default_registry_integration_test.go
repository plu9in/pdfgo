package cli_integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/plu9in/pdfgo/internal/adapters/inbound/cli"
	"github.com/plu9in/pdfgo/internal/adapters/inbound/handlers"
	"github.com/plu9in/pdfgo/internal/application/usecases"
	"github.com/plu9in/pdfgo/internal/domain/document"
)

// --- Double d'intégration: fake saver qui satisfait ports.DocumentSaver ---

type fakeSaver struct {
	called bool
	gotCfg document.DocumentConfig
	ret    string
	err    error
}

func (f *fakeSaver) Save(_ context.Context, cfg document.DocumentConfig) (string, error) {
	f.called = true
	f.gotCfg = cfg
	if f.err != nil {
		return "", f.err
	}
	return f.ret, nil
}

// --- Tests d'intégration légers ---

func Test_CLI_DefaultRegistry_Routes_SaveDocument(t *testing.T) {
	// Arrange: UC réel avec saver fake
	fs := &fakeSaver{ret: "out/test.pdf"}
	uc := usecases.NewSaveDocument(fs)

	// Wiring par défaut: registre avec SaveDocumentHandler
	reg := cli.NewDefaultRegistry(uc)

	// YAML d'entrée pour le CLI
	yamlContent := []byte(`
kind: save_document
spec:
  paper: "A4"
  output:
    format: "pdf"
    path: "out/test.pdf"
`)
	tmpdir := t.TempDir()
	yamlPath := filepath.Join(tmpdir, "req.yaml")
	if err := os.WriteFile(yamlPath, yamlContent, 0o600); err != nil {
		t.Fatalf("write temp yaml: %v", err)
	}

	// Act
	code := cli.Run([]string{"-f", yamlPath}, reg)

	// Assert
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !fs.called {
		t.Fatalf("expected saver to be called via UC")
	}
	if fs.gotCfg.Paper != "A4" || fs.gotCfg.Output.Path != "out/test.pdf" {
		t.Fatalf("unexpected cfg routed to saver: %+v", fs.gotCfg)
	}
}

// Sanity: le handler exposé est bien enregistré
func Test_SaveDocumentHandler_IsRegistered(t *testing.T) {
	fs := &fakeSaver{}
	uc := usecases.NewSaveDocument(fs)
	reg := cli.NewDefaultRegistry(uc)

	h, ok := reg.Get("save_document")
	if !ok {
		t.Fatalf("expected 'save_document' handler to be registered")
	}
	if h.Kind() != "save_document" {
		t.Fatalf("unexpected kind: %s", h.Kind())
	}

	_ = handlers.NewSaveDocumentHandler // empêche l’optimisation d'import
}
