package create

import (
	"exp1/internal/commands/startCmd/interfaces"
	roottimeline "exp1/internal/recorder/root-timeline"
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
}

func NewCreate(event fsnotify.Event, watcher interfaces.IWatcher) *Create{
	return &Create{
		Event: event,
		Watcher: watcher,
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
		Action: "create",
		Timestamp: time.Now(),
	}

	// add file to .rec/root-timeline
	err = roottimeline.Save(data)
	if err != nil{
		log.Fatal("createEvent: ",err)
	}
}