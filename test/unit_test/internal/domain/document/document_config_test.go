package document_test

import (
	"context"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/plu9in/pdfgo/internal/domain/document"
)

func TestDocumentConfig_YAML_RoundTrip_And_SaveCallsSaver(t *testing.T) {
	// YAML de config minimal
	yml := []byte(`
name: "My doc"
paper: "A4"
orientation: "portrait"
margins:
  top: 36
  right: 36
  bottom: 36
  left: 36
metadata:
  author: "Pierre"
  title: "Essai"
output:
  format: "pdf"
  path: "out/test.pdf"
`)

	var cfg document.DocumentConfig
	if err := yaml.Unmarshal(yml, &cfg); err != nil {
		t.Fatalf("yaml unmarshal: %v", err)
	}

	// Fake saver: capture l'appel et renvoie un chemin
	called := false
	var gotCfg document.DocumentConfig
	saver := func(ctx context.Context, c document.DocumentConfig) (string, error) {
		called = true
		gotCfg = c
		return c.Output.Path, nil
	}

	out, err := cfg.Save(context.Background(), saver)
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if !called {
		t.Fatalf("expected saver to be called")
	}
	if out != "out/test.pdf" {
		t.Fatalf("unexpected out path: %s", out)
	}
	if gotCfg.Paper != "A4" || gotCfg.Output.Format != "pdf" {
		t.Fatalf("unexpected cfg passed to saver: %+v", gotCfg)
	}
}

func TestDocumentConfig_ValidateErrors(t *testing.T) {
	// Paper manquant
	cfg := document.DocumentConfig{
		Output: document.Output{Format: "pdf", Path: "out.pdf"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected error on missing paper")
	}

	// Paper inconnu
	cfg = document.DocumentConfig{
		Paper:  "A7", // non supporté
		Output: document.Output{Format: "pdf", Path: "out.pdf"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected error on unknown paper")
	}

	// Format manquant
	cfg = document.DocumentConfig{
		Paper:  "A4",
		Output: document.Output{Path: "out.pdf"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected error on missing output.format")
	}

	// Path manquant
	cfg = document.DocumentConfig{
		Paper:  "A4",
		Output: document.Output{Format: "pdf"},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatalf("expected error on missing output.path")
	}
}
