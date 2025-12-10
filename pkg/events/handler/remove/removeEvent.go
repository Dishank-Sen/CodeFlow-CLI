package remove

import (
	"context"
	roottimeline "exp1/internal/recorder/root-timeline"
	"exp1/internal/types"
	"exp1/pkg/interfaces"
	"exp1/utils/log"
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Remove struct{
    Event fsnotify.Event
	Watcher interfaces.IWatcher
	Ctx context.Context
}

func NewRemove(ctx context.Context, event fsnotify.Event, watcher interfaces.IWatcher) *Remove{
	return &Remove{
		Event: event,
		Watcher: watcher,
		Ctx: ctx,
	}
}

func (r *Remove) Trigger() error{
	path := r.Event.Name
	msg := fmt.Sprintf("file removed: %s", path)
	log.Info(r.Ctx, msg)

	var data = types.FileRecord{
		File: path,
		Action: "remove",
		Timestamp: time.Now(),
	}

	// add file to .rec/root-timeline
	return roottimeline.Save(data)
}