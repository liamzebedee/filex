package bookmarks

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"filex/i18n"
)

// Bookmark represents a sidebar bookmark entry.
type Bookmark struct {
	Name      string
	Path      string
	Icon      string
	UserAdded bool
}

// BookmarkManager manages default and user-added bookmarks.
type BookmarkManager struct {
	defaults []Bookmark
	user     []Bookmark
	filePath string
}

func NewBookmarkManager() *BookmarkManager {
	bm := &BookmarkManager{}
	bm.initDefaults()
	bm.initFilePath()
	bm.loadUser()
	return bm
}

func (bm *BookmarkManager) initDefaults() {
	home, _ := os.UserHomeDir()
	bm.defaults = []Bookmark{
		{Name: i18n.T("Home"), Path: home, Icon: "user-home"},
		{Name: i18n.T("Desktop"), Path: filepath.Join(home, "Desktop"), Icon: "user-desktop"},
		{Name: i18n.T("Documents"), Path: filepath.Join(home, "Documents"), Icon: "folder-documents"},
		{Name: i18n.T("Downloads"), Path: filepath.Join(home, "Downloads"), Icon: "folder-download"},
		{Name: i18n.T("Music"), Path: filepath.Join(home, "Music"), Icon: "folder-music"},
		{Name: i18n.T("Pictures"), Path: filepath.Join(home, "Pictures"), Icon: "folder-pictures"},
		{Name: i18n.T("Videos"), Path: filepath.Join(home, "Videos"), Icon: "folder-videos"},
		{Name: i18n.T("Trash"), Path: filepath.Join(home, ".local/share/Trash/files"), Icon: "user-trash"},
		{Name: i18n.T("File System"), Path: "/", Icon: "drive-harddisk"},
	}

	// Filter to only include directories that exist
	var filtered []Bookmark
	for _, b := range bm.defaults {
		if _, err := os.Stat(b.Path); err == nil {
			filtered = append(filtered, b)
		}
	}
	bm.defaults = filtered
}

func (bm *BookmarkManager) initFilePath() {
	configDir, _ := os.UserConfigDir()
	dir := filepath.Join(configDir, "filex")
	os.MkdirAll(dir, 0755)
	bm.filePath = filepath.Join(dir, "bookmarks.txt")
}

func (bm *BookmarkManager) loadUser() {
	f, err := os.Open(bm.filePath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		path := parts[0]
		name := filepath.Base(path)
		if len(parts) > 1 {
			name = parts[1]
		}
		bm.user = append(bm.user, Bookmark{
			Name:      name,
			Path:      path,
			Icon:      "folder",
			UserAdded: true,
		})
	}
}

func (bm *BookmarkManager) save() {
	f, err := os.Create(bm.filePath)
	if err != nil {
		return
	}
	defer f.Close()

	for _, b := range bm.user {
		fmt.Fprintf(f, "%s|%s\n", b.Path, b.Name)
	}
}

// All returns all bookmarks (defaults + user).
func (bm *BookmarkManager) All() []Bookmark {
	all := make([]Bookmark, 0, len(bm.defaults)+len(bm.user))
	all = append(all, bm.defaults...)
	all = append(all, bm.user...)
	return all
}

// Add adds a user bookmark and persists it.
func (bm *BookmarkManager) Add(b Bookmark) {
	// Don't add duplicates
	for _, existing := range bm.user {
		if existing.Path == b.Path {
			return
		}
	}
	for _, existing := range bm.defaults {
		if existing.Path == b.Path {
			return
		}
	}
	bm.user = append(bm.user, b)
	bm.save()
}

// Remove removes a user bookmark by path.
func (bm *BookmarkManager) Remove(path string) {
	var filtered []Bookmark
	for _, b := range bm.user {
		if b.Path != path {
			filtered = append(filtered, b)
		}
	}
	bm.user = filtered
	bm.save()
}
