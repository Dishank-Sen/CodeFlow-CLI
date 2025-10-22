package events

import (
	"exp1/internal/commands/startCmd/events/handler/create"
	"exp1/internal/commands/startCmd/events/handler/remove"
	"exp1/internal/commands/startCmd/events/handler/rename"
	"exp1/internal/commands/startCmd/events/handler/write"
	"exp1/internal/commands/startCmd/interfaces"
	"exp1/internal/debounce"
	"exp1/internal/types"
	"fmt"
	"time"
	"github.com/fsnotify/fsnotify"
)

type Events struct{
	watcher interfaces.IWatcher
	debouncer *debounce.Debouncer
	State map[string]types.FileRecord
	Unsaved map[string]bool
	RenameFile map[string]time.Time
}

func NewEvents(w interfaces.IWatcher) *Events{
	return &Events{
		watcher: w,
		debouncer: debounce.NewDebouncer(),
		State: make(map[string]types.FileRecord),
		Unsaved: make(map[string]bool),
		RenameFile: make(map[string]time.Time),
	}
}

func (e *Events) Create(event fsnotify.Event) {
	fmt.Println("create event:", event)

	// Check for recent rename events (within 1 second)
	for oldPath, t := range e.RenameFile {
		if time.Since(t) < 1*time.Second {
			fmt.Printf("Detected rename: %s → %s\n", oldPath, event.Name)

			// Clear entry
			delete(e.RenameFile, oldPath)

			newPath := event.Name

			// Call rename handler instead of create
			renameHandler := rename.NewRename(oldPath, newPath, e.watcher)
			renameHandler.RenameTriggered()
			return
		}
	}

	// Normal create event if no recent rename
	createHandler := create.NewCreate(event, e.watcher)
	createHandler.CreateTriggered()
}

func (e *Events) Remove(event fsnotify.Event){
	fmt.Println("remove event:",event)
	remove := remove.NewRemove(event, e.watcher)
	remove.RemoveTriggered()
}

func (e *Events) Rename(event fsnotify.Event){
	fmt.Println("rename event:",event)
	e.RenameFile[event.Name] = time.Now()
}

func (e *Events) Write(event fsnotify.Event){
	path := event.Name

	// Debounce per file path
	e.debouncer.Debounce(path, 2*time.Second, func() {
		writeHandler := write.NewWrite(event, e.watcher, e.State, e.Unsaved)
		writeHandler.WriteTriggered()
	})
}

func (e *Events) Flush(){
	event := fsnotify.Event{}  // empty event
	writeHandler := write.NewWrite(event, e.watcher, e.State, e.Unsaved)
	writeHandler.Flush()
}