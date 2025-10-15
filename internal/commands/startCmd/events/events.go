package events

import (
	"exp1/internal/commands/startCmd/events/handler/chmod"
	"exp1/internal/commands/startCmd/events/handler/create"
	"exp1/internal/commands/startCmd/events/handler/remove"
	"exp1/internal/commands/startCmd/events/handler/rename"
	"exp1/internal/commands/startCmd/events/handler/write"
	"exp1/internal/commands/startCmd/interfaces"
	"exp1/internal/debounce"
	"exp1/internal/types"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Events struct{
	watcher interfaces.IWatcher
	debouncer *debounce.Debouncer
	State map[string]types.FileRecord
	Unsaved map[string]bool
}

func NewEvents(w interfaces.IWatcher) *Events{
	return &Events{
		watcher: w,
		debouncer: debounce.NewDebouncer(),
		State: make(map[string]types.FileRecord),
		Unsaved: make(map[string]bool),
	}
}

func (e *Events) Create(event fsnotify.Event){
	// create object
	create := create.NewCreate(event, e.watcher)

	// call create trigger function
	create.CreateTriggered()
}

func (e *Events) Remove(event fsnotify.Event){
	remove := remove.NewRemove(event, e.watcher)
	remove.RemoveTriggered()
}

func (e *Events) Rename(event fsnotify.Event){
	rename := rename.NewRename(event, e.watcher)
	rename.RenameTriggered()
}

func (e *Events) Write(event fsnotify.Event){
	path := event.Name

	// Debounce per file path
	e.debouncer.Debounce(path, 2*time.Second, func() {
		writeHandler := write.NewWrite(event, e.watcher, e.State, e.Unsaved)
		writeHandler.WriteTriggered()
	})
}

func (e *Events) Chmod(event fsnotify.Event){
	chmod := chmod.NewChmod(event, e.watcher)
	chmod.ChmodTriggered()
}

func (e *Events) Flush(){
	event := fsnotify.Event{}  // empty event
	writeHandler := write.NewWrite(event, e.watcher, e.State, e.Unsaved)
	writeHandler.Flush()
}