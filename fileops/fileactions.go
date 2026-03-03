package fileops

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// OpenFile opens a file with the default application.
func OpenFile(path string) error {
	return exec.Command("xdg-open", path).Start()
}

// OpenTerminal opens a terminal emulator in the given directory.
func OpenTerminal(dir string) error {
	// Try common terminal emulators
	terminals := []string{
		"x-terminal-emulator",
		"gnome-terminal",
		"konsole",
		"xfce4-terminal",
		"xterm",
	}
	for _, term := range terminals {
		if path, err := exec.LookPath(term); err == nil {
			cmd := exec.Command(path)
			cmd.Dir = dir
			return cmd.Start()
		}
	}
	return fmt.Errorf("no terminal emulator found")
}

// Rename renames a file or directory.
func Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// NewFolder creates a new directory.
func NewFolder(path string) error {
	return os.MkdirAll(path, 0755)
}

// TrashFile moves a file to the freedesktop trash.
func TrashFile(path string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	trashDir := filepath.Join(home, ".local", "share", "Trash")
	filesDir := filepath.Join(trashDir, "files")
	infoDir := filepath.Join(trashDir, "info")

	os.MkdirAll(filesDir, 0700)
	os.MkdirAll(infoDir, 0700)

	baseName := filepath.Base(path)
	trashName := baseName
	trashPath := filepath.Join(filesDir, trashName)

	// Handle name collisions in trash
	counter := 1
	for {
		if _, err := os.Stat(trashPath); os.IsNotExist(err) {
			break
		}
		ext := filepath.Ext(baseName)
		name := baseName[:len(baseName)-len(ext)]
		trashName = fmt.Sprintf("%s.%d%s", name, counter, ext)
		trashPath = filepath.Join(filesDir, trashName)
		counter++
	}

	// Write .trashinfo file
	infoContent := fmt.Sprintf("[Trash Info]\nPath=%s\nDeletionDate=%s\n",
		path,
		time.Now().Format("2006-01-02T15:04:05"),
	)
	infoPath := filepath.Join(infoDir, trashName+".trashinfo")
	if err := os.WriteFile(infoPath, []byte(infoContent), 0600); err != nil {
		return err
	}

	// Move file to trash
	if err := os.Rename(path, trashPath); err != nil {
		// Cross-device: copy and delete
		if err := copyFileOrDir(path, trashPath); err != nil {
			os.Remove(infoPath)
			return err
		}
		os.RemoveAll(path)
	}

	return nil
}

// Unzip extracts a zip file to the destination directory.
func Unzip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		// Prevent ZipSlip
		if !filepath.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			// Also allow the destDir itself
			if fpath != filepath.Clean(destDir) {
				return fmt.Errorf("illegal file path: %s", fpath)
			}
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
