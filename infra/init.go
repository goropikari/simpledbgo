package infra

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// InitServer initializes server.
func InitServer(config Config) error {
	if err := os.MkdirAll(config.DBPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(config.DBPath)
	if err != nil {
		return err
	}

	// remove temporary files.
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "temp") {
			if err := os.Remove(filepath.Join(config.DBPath, file.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}
