package watcher

import (
	"bufio"
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

func (w *Watch) Start(){
	w.filterFiles("./")

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

func (w *Watch) eventLoop(){
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Ignore .rec or ignored files
			if w.matchesIgnore(event.Name, w.getIgnoredFiles()) {
				continue
			}

			// Route to correct handler
			if event.Op&fsnotify.Create == fsnotify.Create {
				// fmt.Println("create event triggered")
				if w.events == nil{
					log.Fatal("event is nil")
					break
				}
				w.events.Create(event)
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("write event triggered")
				if w.events == nil{
					log.Fatal("event is nil")
					break
				}
				w.events.Write(event)
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Println("remove event triggered")
				if w.events == nil{
					log.Fatal("event is nil")
					break
				}
				w.events.Remove(event)
			}
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				fmt.Println("rename event triggered")
				if w.events == nil{
					log.Fatal("event is nil")
					break
				}
				w.events.Rename(event)
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Watcher error:", err)
		}
	}
}