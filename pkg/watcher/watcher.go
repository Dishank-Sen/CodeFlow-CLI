package watcher

import (
	"bufio"
	"context"
	"exp1/internal/commands/startCmd/interfaces"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watch struct{
	watcher *fsnotify.Watcher
	events interfaces.IEvents
}

func NewWatcher() *Watch{
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &Watch{
		watcher: watcher,
		events: nil,
	}
}

func (w *Watch) SetEvents(e interfaces.IEvents) {
    w.events = e
}

func (w *Watch) Start() error{
	err := w.filterFiles("./")
	if err != nil{
		return err
	}
	// here code will be blocked
	w.eventLoop()
}

// this removes all files mentioned in .recignore

func (w *Watch) filterFiles(root string) error {
    ignoredPatterns := w.getIgnoredFiles()

    return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // If directory matches ignore pattern, skip it entirely
        if info.IsDir() && w.matchesIgnore(path, ignoredPatterns) {
            return filepath.SkipDir
        }

        // Otherwise, add the directory
        return w.AddDirToWatcher(path, info)
    })
}

func (w *Watch) getIgnoredFiles() []string{
	ignoredPatterns, err := w.loadIgnore(filepath.Join("./", ".recignore"))
	if err != nil && !os.IsNotExist(err) { 
        // ignore error if .recignore not found
		fmt.Println(err)
		return nil
	}
	return ignoredPatterns
}

func (w *Watch) loadIgnore(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns, scanner.Err()
}

func (w *Watch) matchesIgnore(path string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, _ := filepath.Match(pattern, filepath.Base(path))
		if matched {
			return true
		}
		// Handle directory patterns like "vendor/"
		if strings.HasSuffix(pattern, "/") && strings.Contains(path, strings.TrimSuffix(pattern, "/")) {
			return true
		}
	}
	return false
}

// add a dir to be watched

func (w *Watch) AddDirToWatcher(path string, info os.FileInfo) error{
	// Add directories to watcher
	if info.IsDir() {
		fmt.Println("Watching:", path)
		return w.watcher.Add(path)
	}
	return nil
}

// this loop runs forever until termination

func (w *Watch) eventLoop(ctx context.Context) error {
	if w == nil || w.watcher == nil {
		return fmt.Errorf("watcher not initialized")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-w.watcher.Events:
			if !ok {
				// events channel closed by fsnotify; treat as clean shutdown
				return nil
			}

			// ignore logic: replace with your matchesIgnore implementation
			if w.matchesIgnore(event.Name, w.getIgnoredFiles()) {
				continue
			}

			// route to handlers; use != 0 to test bitflags
			if event.Op&fsnotify.Create != 0 {
				if err := w.safeCallCreate(event); err != nil {
					return err
				}
			}
			if event.Op&fsnotify.Write != 0 {
				if err := w.safeCallWrite(event); err != nil {
					return err
				}
			}
			if event.Op&fsnotify.Remove != 0 {
				if err := w.safeCallRemove(event); err != nil {
					return err
				}
			}
			if event.Op&fsnotify.Rename != 0 {
				if err := w.safeCallRename(event); err != nil {
					return err
				}
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				// errors channel closed: treat as clean shutdown
				return nil
			}
			// propagate watcher error to caller
			return fmt.Errorf("fsnotify error: %w", err)
		}
	}
}

// helper wrappers that avoid nil deref and let you convert handler panics to errors if desired
func (w *Watch) safeCallCreate(ev fsnotify.Event) error {
	if w.events == nil {
		// decide policy: return error so caller can log/exit, or ignore
		return fmt.Errorf("events handler is nil")
	}
	w.events.Create(ev)
	return nil
}
func (w *Watch) safeCallWrite(ev fsnotify.Event) error {
	if w.events == nil {
		return fmt.Errorf("events handler is nil")
	}
	w.events.Write(ev)
	return nil
}
func (w *Watch) safeCallRemove(ev fsnotify.Event) error {
	if w.events == nil {
		return fmt.Errorf("events handler is nil")
	}
	w.events.Remove(ev)
	return nil
}
func (w *Watch) safeCallRename(ev fsnotify.Event) error {
	if w.events == nil {
		return fmt.Errorf("events handler is nil")
	}
	w.events.Rename(ev)
	return nil
}

// stubbed functions you referenced — implement them in your package
func (w *Watch) matchesIgnore(name string, ignored []string) bool {
	// TODO: your ignore logic
	return false
}
func (w *Watch) getIgnoredFiles() []string {
	// TODO: return ignored patterns
	return nil
}