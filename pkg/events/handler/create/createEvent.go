package create

import (
	"context"
	"exp1/pkg/interfaces"
	roottimeline "exp1/internal/recorder/root-timeline"
	"exp1/internal/types"
	"exp1/utils/log"
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Create struct{
	Event fsnotify.Event
	Watcher interfaces.IWatcher
	Ctx context.Context
}

func NewCreate(ctx context.Context, event fsnotify.Event, watcher interfaces.IWatcher) *Create{
	return &Create{
		Event: event,
		Watcher: watcher,
		Ctx: ctx,
	}
}

func (c *Create) Trigger() error{
	ctx := c.Ctx

	path := c.Event.Name
	info, err := os.Stat(path)
	if err != nil{
		return err
	}
	if info.IsDir(){
		msg := fmt.Sprintf("folder created: %s", path)
		log.Info(ctx, msg)
		// add folder to watcher
		return c.Watcher.AddDirToWatcher(ctx, path, info)
	}
	msg := fmt.Sprintf("file created: %s", path)
	log.Info(ctx, msg)

	var data = types.FileRecord{
		File: path,
		Action: "create",
		Timestamp: time.Now(),
	}

	// add file to .rec/root-timeline
	return roottimeline.Save(data)
}