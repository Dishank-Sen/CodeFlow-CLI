package write

import (
	"exp1/internal/commands/startCmd/interfaces"
	"exp1/internal/recorder/blob"
	"exp1/internal/recorder/history"
	"exp1/internal/types"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Write struct {
	Event   fsnotify.Event
	Watcher interfaces.IWatcher
	State   map[string]types.FileRecord
	History *history.History
	Blob    *blob.Blob
	Unsaved map[string]bool
}

func NewWrite(event fsnotify.Event, watcher interfaces.IWatcher, state map[string]types.FileRecord, unsaved map[string]bool) *Write {
	return &Write{
		Event:   event,
		Watcher: watcher,
		History: history.NewHistory(),
		State:   state,
		Blob:    blob.NewBlob(),
		Unsaved: unsaved,
	}
}

func (w *Write) WriteTriggered() {
	path := w.Event.Name
	fmt.Println("file write:", path)

	info, err := os.Stat(path)
	if err != nil {
		log.Printf("Failed to stat file (writeEvent.go): %v\n", err)
		return
	}
	size := info.Size()

	// Read new content
	newContentBytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file content: %v\n", err)
		return
	}
	newContent := string(newContentBytes)

	record, exists := w.State[path]

	// ===== Case 1: First time seeing this file =====
	if !exists {
		fmt.Println("🆕 New file detected, creating snapshot record")

		// getting blob path
		blobPath := w.Blob.CreateBlobFromPath(path)

		data := types.FileRecord{
			File:                path,
			Type:                "snapshot",
			Action:              "write",
			PrevSize:            size,
			CurrentSize:         size,
			PreviousFileContent: newContent,
			Timestamp:           time.Now(),
		}
		w.State[path] = data

		// save to history
		historyData := types.FileRecord{
			File: path,
			Type: "snapshot",
			Action: "write",
			IsBlobType: true,
			Blob: blobPath,
			Timestamp: time.Now(),
		}
		err := w.History.Create(path, historyData)

		if err != nil{
			log.Fatal("writeEvent:",err)
		}
		fmt.Println("history created for write snpashot!")
		
		w.Unsaved[path] = false

		return
	}

	// if file already tracked
	record.CurrentSize = size
	fmt.Println("Current size:", record.CurrentSize)
	fmt.Println("Previous size:", record.PrevSize)

	threshold, err := strconv.Atoi(os.Getenv("CODE_THRESHOLD"))
	if err != nil {
		log.Fatal("Invalid CODE_THRESHOLD:", err)
	}

	// Only compute diff if file changed significantly
	if math.Abs(float64(record.CurrentSize)-float64(record.PrevSize)) > float64(threshold) {
		fmt.Println("Significant change detected!")

		currentContentByte, err := os.ReadFile(path)
		if err != nil{
			log.Fatal("error reading file (writeFile.go):",err)
		}

		currentContent := string(currentContentByte)
		previousContent := record.PreviousFileContent

		// create patch
		dmp := diffmatchpatch.New()
		patch := dmp.PatchMake(previousContent, currentContent)

		// save patch as blob
		patchText := dmp.PatchToText(patch)
		blobPath := w.Blob.CreateBlobFromContent(patchText)

		// save history
		historyData := types.FileRecord{
			File: path,
			Type: "delta",
			Action: "write",
			IsBlobType: true,
			Blob: blobPath,
			Timestamp: time.Now(),
		}

		err = w.History.Create(path, historyData)

		if err != nil{
			log.Fatal("writeEvent:",err)
		}
		fmt.Println("history created for write delta!")

		record.CurrentSize = size
		record.PrevSize = size
		record.PreviousFileContent = currentContent
		w.State[path] = record
		w.Unsaved[path] = false

		return
	} else {
		fmt.Println("No significant change detected")
		w.Unsaved[path] = true
		return
	}
}

func (w *Write) Flush(){
	// save snapshot file for every unsaved changes
	var unsavedFiles []string
	for key, value := range w.Unsaved{
		if value{
			unsavedFiles = append(unsavedFiles, key)
		}
	}

	for _, path := range unsavedFiles{
		blobPath := w.Blob.CreateBlobFromPath(path)

		historyData := types.FileRecord{
			File: path,
			Type: "snapshot",
			Action: "write",
			IsBlobType: true,
			Blob: blobPath,
			Timestamp: time.Now(),
		}
		err := w.History.Create(path, historyData)

		if err != nil{
			log.Fatal("error while flushing the file (writeEvent.go):",err)
		}
	}
}