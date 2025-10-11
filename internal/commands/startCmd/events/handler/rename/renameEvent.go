package rename

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

type Rename struct{
	Event fsnotify.Event
	Watcher interfaces.IWatcher
	History *history.History
}

func NewRename(event fsnotify.Event, watcher interfaces.IWatcher) *Rename{
	return &Rename{
		Event: event,
		Watcher: watcher,
		History: history.NewHistory(),
	}
}

func (r *Rename) RenameTriggered(){
	path := r.Event.Name
	info, err := os.Stat(path)
	if err != nil{
		panic(err)
	}
	if info.IsDir(){
		fmt.Println("folder renamed: ",path)
		// add folder to watcher
		r.Watcher.AddDirToWatcher(path, info)
		return
	}
	fmt.Println("file renamed: ",path)

	var data = types.FileRecord{
		File: path,
		Action: "rename",
		IsBlobType: false,
		Timestamp: time.Now(),
	}

	// add file to .rec/history
	err = r.History.Create(path, data)
	if err != nil{
		log.Fatal("removeEvent: ",err)
	}
}
