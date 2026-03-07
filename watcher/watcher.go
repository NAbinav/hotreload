package watcher

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fw   *fsnotify.Watcher
	Root string
}

var ignores = []string{
	".git",
	"node_modules",
	"bin",
	"vendor",
}

func New(root string) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, fmt.Errorf("root directory %q does not exist", root)
	}

	w := &Watcher{
		fw:   fw,
		Root: root,
	}

	err = w.addRecursive(root)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func shouldIgnore(path string) bool {
	base := filepath.Base(path)
	if strings.HasSuffix(base, ".swp") ||
		strings.HasSuffix(base, "~") ||
		strings.HasPrefix(base, "#") {
		return true
	}
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		for _, ig := range ignores {
			if part == ig {
				return true
			}
		}
	}
	return false
}

func (w *Watcher) addRecursive(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {

			if shouldIgnore(path) {
				return filepath.SkipDir
			}

			err := w.fw.Add(path)
			if err != nil {
				return err
			}

			slog.Info("watching", "dir", path)
		}

		return nil
	})
}

func (w *Watcher) Watch(onChange func()) {
	for {
		select {
		case event, ok := <-w.fw.Events:
			if !ok {
				return
			}
			if shouldIgnore(event.Name) {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				slog.Info("file change", "path", event.Name)
				onChange()
			}
			if event.Op&fsnotify.Create != 0 {

				info, err := os.Stat(event.Name)

				if err == nil && info.IsDir() {
					w.addRecursive(event.Name)
				}
			}

		case err := <-w.fw.Errors:
			if strings.Contains(err.Error(), "no such file") {
				continue
			}
			slog.Error("watch error", "err", err)

		}
	}
}

func (w *Watcher) Close() {
	w.fw.Close()
}
