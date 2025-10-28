package main

import (
	"os"

	"github.com/plu9in/pdfgo/internal/adapters/inbound/cli"
)

func main() {
	reg := cli.NewRegistry()
	// Prochaine étape: reg.Register(NewGeneratePDFHandler(...))
	os.Exit(cli.Run(os.Args[1:], reg))
}
