package remove

import (
	"exp1/internal/commands/startCmd/interfaces"
	"exp1/internal/recorder/history"
	roottimeline "exp1/internal/recorder/root-timeline"
	"exp1/internal/types"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Remove struct{
    Event fsnotify.Event
	Watcher interfaces.IWatcher
	History *history.History
}

func NewRemove(event fsnotify.Event, watcher interfaces.IWatcher) *Remove{
	return &Remove{
		Event: event,
		Watcher: watcher,
		History: history.NewHistory(),
	}
}

func (r *Remove) RemoveTriggered(){
	path := r.Event.Name
	info, err := os.Stat(path)
	if err != nil{
		panic(err)
	}
	if info.IsDir(){
		fmt.Println("folder removed: ",path)
		return
	}
	fmt.Println("file removed: ",path)

	var data = types.FileRecord{
		File: path,
		Action: "remove",
		Timestamp: time.Now(),
	}

	// add file to .rec/history
	err = roottimeline.Save(data)
	if err != nil{
		log.Fatal("removeEvent: ",err)
	}
}