package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListFilesBySize(t *testing.T) {
	dir := t.TempDir()
	paths := []struct {
		name string
		size int
	}{
		{"small.txt", 10},
		{"large.txt", 100},
	}
	for _, p := range paths {
		fp := filepath.Join(dir, p.name)
		content := make([]byte, p.size)
		if err := os.WriteFile(fp, content, 0644); err != nil {
			t.Fatalf("write %s: %v", fp, err)
		}
	}

	files, err := listFilesBySize(dir, 0)
	if err != nil {
		t.Fatalf("listFilesBySize returned error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0].Size < files[1].Size {
		t.Fatalf("files not sorted by size")
	}
}
