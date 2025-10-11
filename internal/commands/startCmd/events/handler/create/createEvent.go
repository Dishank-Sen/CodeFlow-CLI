package create

import (
	"exp1/internal/commands/startCmd/interfaces"
	"exp1/internal/recorder/history"
	"exp1/internal/types"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Create struct{
	Event fsnotify.Event
	Watcher interfaces.IWatcher
	History *history.History
}

func NewCreate(event fsnotify.Event, watcher interfaces.IWatcher) *Create{
	return &Create{
		Event: event,
		Watcher: watcher,
		History: history.NewHistory(),
	}
}

func (c *Create) CreateTriggered(){
	path := c.Event.Name
	info, err := os.Stat(path)
	if err != nil{
		panic(err)
	}
	if info.IsDir(){
		fmt.Println("folder created: ",path)
		// add folder to watcher
		c.Watcher.AddDirToWatcher(path, info)
		return
	}
	fmt.Println("file created: ",path)

	var data = types.FileRecord{
		File: path,
		Type: "snapshot",
		Action: "create",
		IsBlobType: false,
		Timestamp: time.Now(),
	}

	// add file to .rec/history
	err = c.History.Create(path, data)
	if err != nil{
		log.Fatal("createEvent: ",err)
	}
}