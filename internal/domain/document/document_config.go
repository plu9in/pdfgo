package document

import (
	"context"
	"fmt"
)

// Orientation simple (on pourra enrichir plus tard)
type Orientation string

const (
	Portrait  Orientation = "portrait"
	Landscape Orientation = "landscape"
)

// Margins exprimées en points PDF (1pt = 1/72 in)
type Margins struct {
	Top    float64 `yaml:"top"`
	Right  float64 `yaml:"right"`
	Bottom float64 `yaml:"bottom"`
	Left   float64 `yaml:"left"`
}

// Output décrit la sortie du rendu (format + destination)
type Output struct {
	Format string `yaml:"format"` // ex: "pdf", "png", "svg"
	Path   string `yaml:"path"`   // ex: "out/doc1.pdf"
}

// DocumentConfig est sérialisable YAML et configure le document à produire.
type DocumentConfig struct {
	Name        string            `yaml:"name"`                  // label métier optionnel
	Paper       string            `yaml:"paper"`                 // ex: "A4", "Letter"
	Orientation Orientation       `yaml:"orientation,omitempty"` // "portrait" (defaut) | "landscape"
	Margins     Margins           `yaml:"margins,omitempty"`
	Metadata    map[string]string `yaml:"metadata,omitempty"` // auteur, titre, etc.
	Output      Output            `yaml:"output"`             // format + chemin cible
}

// SaverFunc est le contrat minimal injecté dans Save.
// (Prochain step: remplacer par un port domain/ports.DocumentSaver)
type SaverFunc func(ctx context.Context, cfg DocumentConfig) (string, error)

// Validate vérifie les paramètres de base (baby step — simple)
func (c DocumentConfig) Validate() error {
	if c.Paper == "" {
		return fmt.Errorf("paper must be set (e.g., A4, Letter)")
	}
	if _, err := GetPaperSize(c.Paper); err != nil {
		return err
	}
	if c.Orientation == "" {
		c.Orientation = Portrait
	}
	switch c.Orientation {
	case Portrait, Landscape:
	default:
		return fmt.Errorf("invalid orientation: %s", c.Orientation)
	}
	if c.Output.Format == "" {
		return fmt.Errorf("output.format must be set (e.g., pdf)")
	}
	if c.Output.Path == "" {
		return fmt.Errorf("output.path must be set")
	}
	return nil
}

// Save applique la configuration en appelant l'impl. injectée.
// Retourne le chemin de sortie effectif (ou un id) selon l’implémentation.
func (c DocumentConfig) Save(ctx context.Context, saver SaverFunc) (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	return saver(ctx, c)
}
