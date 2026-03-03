package fileops

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewFolder(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "newdir")

	err := NewFolder(dir)
	if err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatal("NewFolder should create directory")
	}
	if !info.IsDir() {
		t.Error("NewFolder should create a directory, not a file")
	}
}

func TestNewFolder_Nested(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "a", "b", "c")

	err := NewFolder(dir)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(dir); err != nil {
		t.Fatal("NewFolder should create nested directories")
	}
}

func TestRename(t *testing.T) {
	tmp := t.TempDir()
	old := filepath.Join(tmp, "old.txt")
	newPath := filepath.Join(tmp, "new.txt")

	os.WriteFile(old, []byte("data"), 0644)

	err := Rename(old, newPath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("Old file should not exist after rename")
	}

	data, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatal("New file should exist after rename")
	}
	if string(data) != "data" {
		t.Errorf("Renamed file content = %q, want 'data'", string(data))
	}
}

func TestTrashFile(t *testing.T) {
	tmp := t.TempDir()

	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	testFile := filepath.Join(tmp, "trashme.txt")
	os.WriteFile(testFile, []byte("bye"), 0644)

	err := TrashFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Trashed file should be removed from original location")
	}

	trashFiles := filepath.Join(tmp, ".local", "share", "Trash", "files")
	trashedFile := filepath.Join(trashFiles, "trashme.txt")
	data, err := os.ReadFile(trashedFile)
	if err != nil {
		t.Fatal("File should be in trash/files")
	}
	if string(data) != "bye" {
		t.Errorf("Trashed file content = %q, want 'bye'", string(data))
	}

	trashInfo := filepath.Join(tmp, ".local", "share", "Trash", "info", "trashme.txt.trashinfo")
	infoData, err := os.ReadFile(trashInfo)
	if err != nil {
		t.Fatal("Trashinfo file should exist")
	}
	if !strings.Contains(string(infoData), "[Trash Info]") {
		t.Error("Trashinfo should contain [Trash Info] header")
	}
	if !strings.Contains(string(infoData), testFile) {
		t.Error("Trashinfo should contain original path")
	}
}

func TestTrashFile_NameCollision(t *testing.T) {
	tmp := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", origHome)

	trashFiles := filepath.Join(tmp, ".local", "share", "Trash", "files")
	os.MkdirAll(trashFiles, 0700)

	os.WriteFile(filepath.Join(trashFiles, "dupe.txt"), []byte("old"), 0644)

	testFile := filepath.Join(tmp, "dupe.txt")
	os.WriteFile(testFile, []byte("new"), 0644)

	err := TrashFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(trashFiles, "dupe.1.txt"))
	if err != nil {
		t.Fatal("Collision file should be renamed in trash")
	}
	if string(data) != "new" {
		t.Errorf("Trashed collision file content = %q, want 'new'", string(data))
	}
}

func TestUnzip(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "test.zip")
	destDir := filepath.Join(tmp, "extracted")
	os.Mkdir(destDir, 0755)

	createTestZip(t, zipPath, map[string]string{
		"hello.txt":      "world",
		"sub/nested.txt": "deep",
	})

	err := Unzip(zipPath, destDir)
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(destDir, "hello.txt"))
	if err != nil {
		t.Fatal("hello.txt should be extracted")
	}
	if string(data) != "world" {
		t.Errorf("hello.txt content = %q, want 'world'", string(data))
	}

	data, err = os.ReadFile(filepath.Join(destDir, "sub", "nested.txt"))
	if err != nil {
		t.Fatal("sub/nested.txt should be extracted")
	}
	if string(data) != "deep" {
		t.Errorf("sub/nested.txt content = %q, want 'deep'", string(data))
	}
}

func createTestZip(t *testing.T, zipPath string, files map[string]string) {
	t.Helper()
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(content))
	}
	zw.Close()
}
