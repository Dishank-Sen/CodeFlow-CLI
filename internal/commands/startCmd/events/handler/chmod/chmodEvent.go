package chmod

import (
	"exp1/internal/commands/startCmd/interfaces"

	"github.com/fsnotify/fsnotify"
)

type Chmod struct{
	Event fsnotify.Event
	Watcher interfaces.IWatcher
}

func NewChmod(event fsnotify.Event, watcher interfaces.IWatcher) *Chmod{
	return &Chmod{
		Event: event,
		Watcher: watcher,
	}
}

func (c *Chmod) ChmodTriggered(){
	
}