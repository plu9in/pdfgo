package handlers_test

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/plu9in/pdfgo/internal/adapters/inbound/handlers"
	"github.com/plu9in/pdfgo/internal/domain/document"
)

type fakeSaveDocumentUC struct {
	called bool
	gotCfg document.DocumentConfig
	ret    string
	err    error
}

func (f *fakeSaveDocumentUC) Execute(ctx any, cfg document.DocumentConfig) (string, error) {
	f.called = true
	f.gotCfg = cfg
	return f.ret, f.err
}

func TestSaveDocumentHandler_Kind(t *testing.T) {
	uc := &fakeSaveDocumentUC{}
	h := handlers.NewSaveDocumentHandler(uc)

	if h.Kind() != "save_document" {
		t.Fatalf("expected kind 'save_document', got %q", h.Kind())
	}
}

func TestSaveDocumentHandler_ParsesYAML_AndCallsUseCase(t *testing.T) {
	uc := &fakeSaveDocumentUC{ret: "out/test.pdf"}
	h := handlers.NewSaveDocumentHandler(uc)

	// YAML d'entrée attendu par le handler: {kind, spec{...DocumentConfig...}}
	yml := []byte(`
kind: save_document
spec:
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
  output:
    format: "pdf"
    path: "out/test.pdf"
`)

	// par sécurité, on valide qu'il est bien du YAML
	var tmp any
	if err := yaml.Unmarshal(yml, &tmp); err != nil {
		t.Fatalf("test yaml is invalid: %v", err)
	}

	if err := h.Handle(yml); err != nil {
		t.Fatalf("Handle returned error: %v", err)
	}
	if !uc.called {
		t.Fatalf("expected use case to be called")
	}
	if uc.gotCfg.Paper != "A4" || uc.gotCfg.Output.Format != "pdf" || uc.gotCfg.Output.Path != "out/test.pdf" {
		t.Fatalf("unexpected cfg passed to UC: %+v", uc.gotCfg)
	}
}

func TestSaveDocumentHandler_BadYAML_ReturnsError(t *testing.T) {
	uc := &fakeSaveDocumentUC{}
	h := handlers.NewSaveDocumentHandler(uc)

	// Pas de bloc spec
	yml := []byte(`kind: save_document`)

	if err := h.Handle(yml); err == nil {
		t.Fatalf("expected error on missing spec")
	}
	if uc.called {
		t.Fatalf("use case should not have been called on bad yaml")
	}
}

func TestSaveDocumentHandler_InvalidConfig_PreventsUseCase(t *testing.T) {
	uc := &fakeSaveDocumentUC{}
	h := handlers.NewSaveDocumentHandler(uc)

	// Paper inconnu => la validation interne doit échouer avant UC
	yml := []byte(`
kind: save_document
spec:
  paper: "A7"
  output:
    format: "pdf"
    path: "out/bad.pdf"
`)

	if err := h.Handle(yml); err == nil {
		t.Fatalf("expected validation error for unknown paper")
	}
	if uc.called {
		t.Fatalf("use case should not have been called on invalid config")
	}
}
