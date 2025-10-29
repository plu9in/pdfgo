package document

import "fmt"

// PaperSize représente une taille de papier (largeur × hauteur en points PDF).
// 1 point = 1/72 inch.
type PaperSize struct {
	Name   string
	Width  float64 // points
	Height float64 // points
}

// PaperSizes — quelques formats usuels (points).
var PaperSizes = map[string]PaperSize{
	"A0":      {"A0", 2383.94, 3370.39},
	"A1":      {"A1", 1683.78, 2383.94},
	"A2":      {"A2", 1190.55, 1683.78},
	"A3":      {"A3", 841.89, 1190.55},
	"A4":      {"A4", 595.28, 841.89},
	"A5":      {"A5", 419.53, 595.28},
	"LETTER":  {"LETTER", 612.0, 792.0},
	"LEGAL":   {"LEGAL", 612.0, 1008.0},
	"TABLOID": {"TABLOID", 792.0, 1224.0},
}

// GetPaperSize renvoie la taille demandée, insensible à la casse.
func GetPaperSize(name string) (PaperSize, error) {
	if name == "" {
		return PaperSize{}, fmt.Errorf("paper size name is empty")
	}
	k := toUpperASCII(name)
	if ps, ok := PaperSizes[k]; ok {
		return ps, nil
	}
	return PaperSize{}, fmt.Errorf("unknown paper size: %s", name)
}

// toUpperASCII : évite une dépendance unicode pour ce besoin simple.
func toUpperASCII(s string) string {
	b := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		b[i] = c
	}
	return string(b)
}
