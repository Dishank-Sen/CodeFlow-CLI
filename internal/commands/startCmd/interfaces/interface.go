package interfaces

import (
	"os"

	"github.com/fsnotify/fsnotify"
)

type IEvents interface{
	Create(event fsnotify.Event)
	Remove(event fsnotify.Event)
	Rename(event fsnotify.Event)
	Write(event fsnotify.Event)
	Chmod(event fsnotify.Event)
}

type IWatcher interface{
	AddDirToWatcher(path string, info os.FileInfo) error
}
