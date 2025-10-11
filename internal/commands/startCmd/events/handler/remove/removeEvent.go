package remove

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
		IsBlobType: false,
		Timestamp: time.Now(),
	}

	// add file to .rec/history
	err = r.History.Create(path, data)
	if err != nil{
		log.Fatal("removeEvent: ",err)
	}
}

// func (h *Handler) Remove(event fsnotify.Event) {
//     fmt.Println("remove event triggered")

//     path := event.Name

//     // Check if we had it in state (to know if it was file or dir)
//     if _, exists := h.State[path]; exists {
//         fmt.Println("file deleted:", path)
//         delete(h.State, path)

//         // record in history
//         err := h.Manager.FileRemoveSnapshot(path)
//         if err != nil {
//             fmt.Println("error recording remove snapshot:", err)
//         }
//     } else {
//         fmt.Println("removed unknown path (maybe dir):", path)
//     }
// }
