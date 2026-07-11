package cleanup

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSweep_RemovesOnlyOldFiles(t *testing.T) {
	dir := t.TempDir()

	// Old file: backdate its modification time well past the cutoff.
	old := filepath.Join(dir, "old.bin")
	if err := os.WriteFile(old, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	backdated := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(old, backdated, backdated); err != nil {
		t.Fatal(err)
	}

	// Fresh file: must survive.
	fresh := filepath.Join(dir, "fresh.bin")
	if err := os.WriteFile(fresh, []byte("fresh"), 0644); err != nil {
		t.Fatal(err)
	}

	// A subdirectory must be left untouched.
	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}

	if removed := Sweep([]string{dir}, time.Hour); removed != 1 {
		t.Fatalf("expected 1 file removed, got %d", removed)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Errorf("old file should have been removed (stat err = %v)", err)
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Errorf("fresh file should remain: %v", err)
	}
	if _, err := os.Stat(sub); err != nil {
		t.Errorf("subdir should remain: %v", err)
	}
}

func TestSweep_MissingDirIsFine(t *testing.T) {
	if removed := Sweep([]string{"/no/such/dir/xyz"}, time.Hour); removed != 0 {
		t.Fatalf("expected 0 removed for missing dir, got %d", removed)
	}
}
