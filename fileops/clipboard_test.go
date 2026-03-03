package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPasteFiles_Copy(t *testing.T) {
	tmp := t.TempDir()
	srcDir := filepath.Join(tmp, "src")
	dstDir := filepath.Join(tmp, "dst")
	os.Mkdir(srcDir, 0755)
	os.Mkdir(dstDir, 0755)

	// Create source file
	srcFile := filepath.Join(srcDir, "test.txt")
	os.WriteFile(srcFile, []byte("content"), 0644)

	err := PasteFiles([]string{srcFile}, dstDir, false)
	if err != nil {
		t.Fatal(err)
	}

	// Source should still exist (copy, not cut)
	if _, err := os.Stat(srcFile); err != nil {
		t.Error("Source file should still exist after copy")
	}

	// Destination should exist
	dstFile := filepath.Join(dstDir, "test.txt")
	data, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatal("Destination file should exist after copy")
	}
	if string(data) != "content" {
		t.Errorf("Copied file content = %q, want 'content'", string(data))
	}
}

func TestPasteFiles_Cut(t *testing.T) {
	tmp := t.TempDir()
	srcDir := filepath.Join(tmp, "src")
	dstDir := filepath.Join(tmp, "dst")
	os.Mkdir(srcDir, 0755)
	os.Mkdir(dstDir, 0755)

	srcFile := filepath.Join(srcDir, "test.txt")
	os.WriteFile(srcFile, []byte("moved"), 0644)

	err := PasteFiles([]string{srcFile}, dstDir, true)
	if err != nil {
		t.Fatal(err)
	}

	// Source should be gone (cut)
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Error("Source file should not exist after cut")
	}

	// Destination should exist
	dstFile := filepath.Join(dstDir, "test.txt")
	data, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatal("Destination file should exist after cut")
	}
	if string(data) != "moved" {
		t.Errorf("Moved file content = %q, want 'moved'", string(data))
	}
}

func TestPasteFiles_Collision(t *testing.T) {
	tmp := t.TempDir()
	srcDir := filepath.Join(tmp, "src")
	dstDir := filepath.Join(tmp, "dst")
	os.Mkdir(srcDir, 0755)
	os.Mkdir(dstDir, 0755)

	// Create source and existing dest with same name
	srcFile := filepath.Join(srcDir, "test.txt")
	os.WriteFile(srcFile, []byte("new"), 0644)
	os.WriteFile(filepath.Join(dstDir, "test.txt"), []byte("existing"), 0644)

	err := PasteFiles([]string{srcFile}, dstDir, false)
	if err != nil {
		t.Fatal(err)
	}

	// Original should be untouched
	data, _ := os.ReadFile(filepath.Join(dstDir, "test.txt"))
	if string(data) != "existing" {
		t.Error("Original file should be untouched")
	}

	// Copy should exist with "(copy)" suffix
	copyFile := filepath.Join(dstDir, "test (copy).txt")
	data, err = os.ReadFile(copyFile)
	if err != nil {
		t.Fatal("Copy file should exist with (copy) suffix")
	}
	if string(data) != "new" {
		t.Errorf("Copy file content = %q, want 'new'", string(data))
	}
}

func TestPasteFiles_CopyDir(t *testing.T) {
	tmp := t.TempDir()
	srcDir := filepath.Join(tmp, "src", "mydir")
	dstDir := filepath.Join(tmp, "dst")
	os.MkdirAll(srcDir, 0755)
	os.Mkdir(dstDir, 0755)

	// Create file inside source dir
	os.WriteFile(filepath.Join(srcDir, "inner.txt"), []byte("inside"), 0644)

	err := PasteFiles([]string{srcDir}, dstDir, false)
	if err != nil {
		t.Fatal(err)
	}

	// Check destination
	innerFile := filepath.Join(dstDir, "mydir", "inner.txt")
	data, err := os.ReadFile(innerFile)
	if err != nil {
		t.Fatal("Inner file should exist in copied directory")
	}
	if string(data) != "inside" {
		t.Errorf("Inner file content = %q, want 'inside'", string(data))
	}
}

func TestUniquePath(t *testing.T) {
	tmp := t.TempDir()

	// Non-existent path returns as-is
	p := filepath.Join(tmp, "new.txt")
	got := uniquePath(p)
	if got != p {
		t.Errorf("uniquePath for non-existent = %q, want %q", got, p)
	}

	// Existing file gets "(copy)"
	os.WriteFile(p, []byte("x"), 0644)
	got = uniquePath(p)
	want := filepath.Join(tmp, "new (copy).txt")
	if got != want {
		t.Errorf("uniquePath for existing = %q, want %q", got, want)
	}

	// If "(copy)" also exists, get "(copy 2)"
	os.WriteFile(want, []byte("x"), 0644)
	got = uniquePath(p)
	want2 := filepath.Join(tmp, "new (copy 2).txt")
	if got != want2 {
		t.Errorf("uniquePath double collision = %q, want %q", got, want2)
	}
}
