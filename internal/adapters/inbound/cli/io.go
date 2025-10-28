package cli

import (
	"io"
	"os"
)

type requestHeader struct {
	Kind string `yaml:"kind"`
}

func readAll(path string) ([]byte, error) {
	var rd io.Reader
	if path == "-" {
		rd = os.Stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		rd = f
	}
	return io.ReadAll(rd)
}
