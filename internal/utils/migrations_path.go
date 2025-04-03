package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetMigrationsPath() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working dir: %v", err)
	}

	rootMarker := "collecting-metrics-alert-service"
	idx := strings.LastIndex(wd, rootMarker)
	if idx == -1 {
		log.Fatalf("project root marker '%s' not found in: %s", rootMarker, wd)
	}

	rootPath := wd[:idx+len(rootMarker)]
	migrationsPath := filepath.Join(rootPath, "internal", "schema")

	return migrationsPath
}
