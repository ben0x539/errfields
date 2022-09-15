package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/ben0x539/errfields"
	"github.com/ben0x539/errfields/zapfields"
)

func main() {
	logger := zap.NewExample()

	if err := CatAllFiles(os.Args[1:]); err != nil {
		zapfields.With(logger, err).Fatal("request failed")
	}
}

func CatAllFiles(paths []string) error {
	for _, path := range paths {
		s, err := os.ReadFile(path)
		if err != nil {
			return errfields.Add(err, "path", path)
		}

		os.Stdout.Write(s)
	}

	return nil
}
