package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var invalidNameCharacters = regexp.MustCompile(`[^a-z0-9]+`)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Migration name: ")

	name, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		if len(name) == 0 {
			fmt.Fprintln(os.Stderr, "read migration name:", err)
			os.Exit(1)
		}
	}

	files, err := createMigration("migrations", name, time.Now())
	if err != nil {
		fmt.Fprintln(os.Stderr, "create migration:", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s and %s\n", files.up, files.down)
}

type migrationFiles struct {
	up   string
	down string
}

func createMigration(directory, rawName string, now time.Time) (migrationFiles, error) {
	name := normalizeName(rawName)
	if name == "" {
		return migrationFiles{}, errors.New("name must contain letters or numbers")
	}

	if err := os.MkdirAll(directory, 0o755); err != nil {
		return migrationFiles{}, fmt.Errorf("create migrations directory: %w", err)
	}

	prefix := now.Format("20060102150405")
	base := prefix + "_" + name
	files := migrationFiles{
		up:   filepath.Join(directory, base+".up.sql"),
		down: filepath.Join(directory, base+".down.sql"),
	}

	for _, path := range []string{files.up, files.down} {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				return migrationFiles{}, fmt.Errorf("migration already exists: %s", path)
			}
			return migrationFiles{}, fmt.Errorf("create %s: %w", path, err)
		}
		if err := file.Close(); err != nil {
			return migrationFiles{}, fmt.Errorf("close %s: %w", path, err)
		}
	}

	return files, nil
}

func normalizeName(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = invalidNameCharacters.ReplaceAllString(value, "_")
	return strings.Trim(value, "_")
}
