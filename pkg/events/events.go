package events

import (
	"context"
	"exp1/pkg/events/handler/rename"
	"exp1/pkg/events/handler/write"
	"exp1/pkg/interfaces"
	"exp1/internal/debounce"
	"exp1/internal/types"
	"exp1/pkg/events/handler/create"
	"exp1/pkg/events/handler/remove"
	"exp1/utils/log"
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
	Ctx context.Context
}

func NewEvents(w interfaces.IWatcher, ctx context.Context) *Events{
	return &Events{
		watcher: w,
		debouncer: debounce.NewDebouncer(),
		State: make(map[string]types.FileRecord),
		Unsaved: make(map[string]bool),
		RenameFile: make(map[string]time.Time),
		Ctx: ctx,
	}
}

func (e *Events) Create(event fsnotify.Event) error{
	// fmt.Println("create event:", event)

	// Check for recent rename events (within 1 second)
	for oldPath, t := range e.RenameFile {
		if time.Since(t) < 1*time.Second {
			msg := fmt.Sprintf("Detected rename: %s → %s", oldPath, event.Name)
			log.Info(e.Ctx, msg)

			// Clear entry
			delete(e.RenameFile, oldPath)

			newPath := event.Name

			// Call rename handler instead of create
			renameHandler := rename.NewRename(e.Ctx, oldPath, newPath, e.watcher)
			renameHandler.Trigger()
			return nil
		}
	}

	// Normal create event if no recent rename
	createHandler := create.NewCreate(e.Ctx, event, e.watcher)
	return createHandler.Trigger()
}

func (e *Events) Remove(event fsnotify.Event) error{
	// fmt.Println("remove event:",event)
	remove := remove.NewRemove(e.Ctx, event, e.watcher)
	return remove.Trigger()
}

func (e *Events) Rename(event fsnotify.Event) error{
	// fmt.Println("rename event:",event)
	e.RenameFile[event.Name] = time.Now()
	return nil
}

func (e *Events) Write(event fsnotify.Event) error{
	path := event.Name
	var err error
	// Debounce per file path
	debounceTime, err := debounce.GetDebounceTime()
	if err != nil{
		return err
	}
	if debounceTime == 0{
		log.Info(e.Ctx, "no debounce time set")
		debounce.SetDebounceTime(2)
		debounceTime, err = debounce.GetDebounceTime()
		if err != nil{
			return err
		}
	}

	e.debouncer.Debounce(path, time.Duration(debounceTime)*time.Second, func() {
		writeHandler := write.NewWrite(e.Ctx, event, e.watcher, e.State, e.Unsaved)
		err = writeHandler.Trigger()
	})
	return err
}

func (e *Events) Flush() error{
	event := fsnotify.Event{}  // empty event
	writeHandler := write.NewWrite(e.Ctx, event, e.watcher, e.State, e.Unsaved)
	return writeHandler.Flush()
}