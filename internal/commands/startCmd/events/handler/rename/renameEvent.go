package rename

import (
	"exp1/internal/commands/startCmd/interfaces"
	"exp1/internal/recorder/history"
	roottimeline "exp1/internal/recorder/root-timeline"
	"exp1/internal/types"
	"log"
	"time"
)

type Rename struct{
	OldPath string
	NewPath string
	Watcher interfaces.IWatcher
	History *history.History
}

func NewRename(oldPath string, newpath string, watcher interfaces.IWatcher) *Rename{
	return &Rename{
		OldPath: oldPath,
		NewPath: newpath,
		Watcher: watcher,
		History: history.NewHistory(),
	}
}

func (r *Rename) RenameTriggered(){

	var data = types.FileRecord{
		File: r.NewPath,
		Action: "rename",
		NewPath: r.NewPath,
		OldPath: r.OldPath,
		Timestamp: time.Now(),
	}

	// add file to .rec/history
	err := roottimeline.Save(data)
	if err != nil{
		log.Fatal("removeEvent: ",err)
	}
}
