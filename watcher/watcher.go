package watcher

import (
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
}

func New(root string) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
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
			slog.Error("watch error", "err", err)

		}
	}
}

func (w *Watcher) Close() {
	w.fw.Close()
}
