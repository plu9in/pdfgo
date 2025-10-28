package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/plu9in/pdfgo/internal/adapters/inbound/cli"
)

type fakeHandler struct {
	kind    string
	called  bool
	payload []byte
	err     error
}

func (f *fakeHandler) Kind() string { return f.kind }
func (f *fakeHandler) Handle(b []byte) error {
	f.called = true
	f.payload = append([]byte(nil), b...)
	return f.err
}

func TestCLI_RoutesToRegisteredHandler_ByKind(t *testing.T) {
	yamlContent := []byte(`kind: generate_pdf
spec:
  documentID: "doc-42"
  template: "simple.tpl"
`)
	tmpdir := t.TempDir()
	yamlPath := filepath.Join(tmpdir, "req.yaml")
	if err := os.WriteFile(yamlPath, yamlContent, 0o600); err != nil {
		t.Fatalf("write temp yaml: %v", err)
	}

	reg := cli.NewRegistry()
	h := &fakeHandler{kind: "generate_pdf"}
	reg.Register(h)

	code := cli.Run([]string{"-f", yamlPath}, reg)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !h.called || len(h.payload) == 0 {
		t.Fatalf("expected handler to be called with payload")
	}
}

func TestCLI_Errors_WhenNoHandlerForKind(t *testing.T) {
	yamlContent := []byte(`kind: unknown_kind`)
	tmpdir := t.TempDir()
	yamlPath := filepath.Join(tmpdir, "req.yaml")
	_ = os.WriteFile(yamlPath, yamlContent, 0o600)

	reg := cli.NewRegistry()
	code := cli.Run([]string{"-f", yamlPath}, reg)
	if code == 0 {
		t.Fatalf("expected non-zero exit code when handler missing")
	}
}

func TestCLI_MissingFlag_ReturnsUsageError(t *testing.T) {
	reg := cli.NewRegistry()
	code := cli.Run([]string{}, reg)
	if code == 0 {
		t.Fatalf("expected usage error due to missing -f")
	}
}
