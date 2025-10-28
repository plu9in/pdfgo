package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func Run(args []string, reg *Registry) int {
	fs := flag.NewFlagSet("pdfgo", flag.ContinueOnError)
	var filePath string
	fs.StringVar(&filePath, "f", "", "Chemin du fichier YAML (ou '-' pour stdin)")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: pdfgo -f <fichier.yaml|->")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "erreur: %v\n", err)
		return 2
	}
	if filePath == "" {
		fmt.Fprintln(os.Stderr, "erreur: option -f requise (chemin YAML ou '-')")
		return 2
	}

	raw, err := readAll(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "erreur de lecture YAML: %v\n", err)
		return 1
	}
	if len(raw) == 0 {
		fmt.Fprintln(os.Stderr, "erreur: YAML vide")
		return 1
	}

	var hdr requestHeader
	if err := yaml.Unmarshal(raw, &hdr); err != nil {
		fmt.Fprintf(os.Stderr, "erreur YAML (entête): %v\n", err)
		return 1
	}
	if hdr.Kind == "" {
		fmt.Fprintln(os.Stderr, "erreur: champ 'kind' manquant dans le YAML")
		return 1
	}

	h, ok := reg.Get(hdr.Kind)
	if !ok {
		fmt.Fprintf(os.Stderr, "erreur: aucun handler enregistré pour kind=%q\n", hdr.Kind)
		return 3
	}
	if err := h.Handle(raw); err != nil {
		if !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "échec handler %q: %v\n", hdr.Kind, err)
		}
		return 1
	}
	return 0
}
