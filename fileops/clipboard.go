package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PasteFiles copies (or moves, when cut) sources into destDir. Pasting a
// directory into itself or its own subtree is skipped — copying a tree
// into one of its descendants would recurse forever.
func PasteFiles(sources []string, destDir string, cut bool) error {
	sep := string(os.PathSeparator)
	for _, src := range sources {
		if src == destDir || strings.HasPrefix(destDir+sep, src+sep) {
			continue
		}

		// Handle name collisions
		dest := uniquePath(filepath.Join(destDir, filepath.Base(src)))

		if cut {
			if err := os.Rename(src, dest); err != nil {
				// Cross-device move: copy then delete
				if err := copyFileOrDir(src, dest); err != nil {
					return err
				}
				os.RemoveAll(src)
			}
		} else {
			if err := copyFileOrDir(src, dest); err != nil {
				return err
			}
		}
	}
	return nil
}

// uniquePath appends "(copy)" or "(copy N)" to avoid collisions.
func uniquePath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	dir := filepath.Dir(path)
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(filepath.Base(path), ext)

	newPath := filepath.Join(dir, name+" (copy)"+ext)
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return newPath
	}

	for i := 2; ; i++ {
		newPath = filepath.Join(dir, name+fmt.Sprintf(" (copy %d)", i)+ext)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath
		}
	}
}

func copyFileOrDir(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyDir(src, dest)
	}
	return copyFile(src, dest)
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	// Preserve permissions
	info, err := os.Stat(src)
	if err == nil {
		os.Chmod(dest, info.Mode())
	}

	return nil
}

func copyDir(src, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dest, info.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dest, entry.Name())
		if err := copyFileOrDir(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}
