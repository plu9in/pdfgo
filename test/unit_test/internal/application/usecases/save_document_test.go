package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/plu9in/pdfgo/internal/application/usecases"
	"github.com/plu9in/pdfgo/internal/domain/document"
)

type fakeSaver struct {
	called bool
	gotCfg document.DocumentConfig
	ret    string
	err    error
}

func (f *fakeSaver) Save(ctx context.Context, cfg document.DocumentConfig) (string, error) {
	f.called = true
	f.gotCfg = cfg
	if f.err != nil {
		return "", f.err
	}
	return f.ret, nil
}

func TestSaveDocument_DelegatesToSaver_AndReturnsPath(t *testing.T) {
	// config valide
	cfg := document.DocumentConfig{
		Name:    "MyDoc",
		Paper:   "A4",
		Margins: document.Margins{Top: 36, Right: 36, Bottom: 36, Left: 36},
		Output:  document.Output{Format: "pdf", Path: "out/ok.pdf"},
	}

	fs := &fakeSaver{ret: "out/ok.pdf"}

	uc := usecases.NewSaveDocument(fs)
	out, err := uc.Execute(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fs.called {
		t.Fatalf("expected saver to be called")
	}
	if out != "out/ok.pdf" {
		t.Fatalf("unexpected returned path: %s", out)
	}
	if fs.gotCfg.Paper != "A4" || fs.gotCfg.Output.Format != "pdf" {
		t.Fatalf("unexpected cfg passed to saver: %+v", fs.gotCfg)
	}
}

func TestSaveDocument_InvalidConfig_PreventsSaverCall(t *testing.T) {
	// paper inconnu => Validate doit échouer et ne pas appeler le saver
	cfg := document.DocumentConfig{
		Paper:  "A7", // non supporté
		Output: document.Output{Format: "pdf", Path: "out/bad.pdf"},
	}

	fs := &fakeSaver{ret: "out/bad.pdf"}

	uc := usecases.NewSaveDocument(fs)
	out, err := uc.Execute(context.Background(), cfg)
	if err == nil {
		t.Fatalf("expected error for invalid config, got nil (out=%q)", out)
	}
	if fs.called {
		t.Fatalf("saver should NOT have been called on invalid config")
	}
}

func TestSaveDocument_PropagatesSaverError(t *testing.T) {
	cfg := document.DocumentConfig{
		Paper:  "A4",
		Output: document.Output{Format: "pdf", Path: "out/err.pdf"},
	}
	fs := &fakeSaver{err: errors.New("disk full")}

	uc := usecases.NewSaveDocument(fs)
	_, err := uc.Execute(context.Background(), cfg)
	if err == nil {
		t.Fatalf("expected saver error to propagate")
	}
}
