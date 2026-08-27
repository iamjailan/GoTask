package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrationDefinitions(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{
		"000010_add_tasks.up.sql",
		"000002_create_customers.up.sql",
		"000010_add_tasks.down.sql",
		"not-a-migration.up.sql",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	definitions, err := migrationDefinitions(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(definitions) != 2 {
		t.Fatalf("migration definitions = %d, want 2", len(definitions))
	}
	if definitions[0].version != 2 || definitions[0].name != "create_customers" {
		t.Errorf("first definition = %#v, want version 2 named create_customers", definitions[0])
	}
	if definitions[1].version != 10 || definitions[1].name != "add_tasks" {
		t.Errorf("second definition = %#v, want version 10 named add_tasks", definitions[1])
	}
}
