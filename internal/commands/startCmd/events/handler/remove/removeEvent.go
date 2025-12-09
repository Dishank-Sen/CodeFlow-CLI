package remove

import (
	"exp1/internal/commands/startCmd/interfaces"
	roottimeline "exp1/internal/recorder/root-timeline"
	"exp1/internal/types"
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Remove struct{
    Event fsnotify.Event
	Watcher interfaces.IWatcher
}

func NewRemove(event fsnotify.Event, watcher interfaces.IWatcher) *Remove{
	return &Remove{
		Event: event,
		Watcher: watcher,
	}
}

func (r *Remove) RemoveTriggered(){
	path := r.Event.Name
	fmt.Println("file removed: ",path)

	var data = types.FileRecord{
		File: path,
		Action: "remove",
		Timestamp: time.Now(),
	}

	// add file to .rec/root-timeline
	err := roottimeline.Save(data)
	if err != nil{
		log.Fatal("removeEvent: ",err)
	}
}