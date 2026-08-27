package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNormalizeName(t *testing.T) {
	got := normalizeName("  Add due-date to Tasks!  ")
	if got != "add_due_date_to_tasks" {
		t.Fatalf("normalizeName() = %q, want %q", got, "add_due_date_to_tasks")
	}
}

func TestCreateMigration(t *testing.T) {
	directory := t.TempDir()
	now := time.Date(2026, time.August, 27, 14, 30, 0, 0, time.UTC)

	files, err := createMigration(directory, "Add due date to tasks", now)
	if err != nil {
		t.Fatalf("createMigration() error = %v", err)
	}

	for _, path := range []string{files.up, files.down} {
		if filepath.Dir(path) != directory {
			t.Errorf("migration directory = %q, want %q", filepath.Dir(path), directory)
		}
		if _, err := os.Stat(path); err != nil {
			t.Errorf("migration file %q was not created: %v", path, err)
		}
	}

	if filepath.Base(files.up) != "20260827143000_add_due_date_to_tasks.up.sql" {
		t.Errorf("up migration = %q", filepath.Base(files.up))
	}
	if filepath.Base(files.down) != "20260827143000_add_due_date_to_tasks.down.sql" {
		t.Errorf("down migration = %q", filepath.Base(files.down))
	}
}
