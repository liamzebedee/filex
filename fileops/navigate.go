// Package fileops contains the filesystem effects of the file manager:
// listing, copying, moving, trashing, extracting. It depends on core for
// the data types but never on the UI.
package fileops

import (
	"os"
	"path/filepath"

	"filex/core"
)

// ListDirectory reads a directory and returns every entry, hidden files
// included — what is shown is decided later by the pure core.Visible
// selector, so toggling filters never touches the disk again.
func ListDirectory(dirPath string) ([]core.FileEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	result := make([]core.FileEntry, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		path := filepath.Join(dirPath, e.Name())
		isDir := info.IsDir()
		if info.Mode()&os.ModeSymlink != 0 {
			// A symlink to a directory should behave like a directory
			// (navigate into it rather than "open" it).
			if target, err := os.Stat(path); err == nil {
				isDir = target.IsDir()
			}
		}
		result = append(result, core.FileEntry{
			Name:    e.Name(),
			Path:    path,
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
			IsDir:   isDir,
			Mode:    info.Mode(),
		})
	}
	return result, nil
}
