package write

import (
	"encoding/json"
	"exp1/internal/commands/startCmd/interfaces"
	diffalgo "exp1/internal/diffAlgo"
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
)

type Write struct {
	Event   fsnotify.Event
	Watcher interfaces.IWatcher
	State   map[string]types.FileRecord
	History *history.History
	Blob    *blob.Blob
}

func NewWrite(event fsnotify.Event, watcher interfaces.IWatcher, state map[string]types.FileRecord) *Write {
	return &Write{
		Event:   event,
		Watcher: watcher,
		History: history.NewHistory(),
		State:   state,
		Blob:    blob.NewBlob(),
	}
}

func (w *Write) WriteTriggered() {
	path := w.Event.Name
	fmt.Println("file write:", path)

	info, err := os.Stat(path)
	if err != nil {
		log.Printf("⚠️ Failed to stat file: %v\n", err)
		return
	}
	size := info.Size()

	// Read new content
	newContentBytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("⚠️ Failed to read file content: %v\n", err)
		return
	}
	newContent := string(newContentBytes)

	record, exists := w.State[path]

	// ===== Case 1: First time seeing this file =====
	if !exists {
		fmt.Println("🆕 New file detected, creating snapshot record")

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
			Timestamp: time.Now(),
		}
		err := w.History.Create(path, historyData)

		if err != nil{
			log.Fatal("writeEvent:",err)
		}
		fmt.Println("history created for write snpashot!")

		return
	}

	// ===== Case 2: File already tracked =====
	record.CurrentSize = size
	fmt.Println("📏 Current size:", record.CurrentSize)
	fmt.Println("📏 Previous size:", record.PrevSize)

	threshold, err := strconv.Atoi(os.Getenv("CODE_THRESHOLD"))
	if err != nil {
		log.Fatal("❌ Invalid CODE_THRESHOLD:", err)
	}

	// Only compute diff if file changed significantly
	if math.Abs(float64(record.CurrentSize)-float64(record.PrevSize)) > float64(threshold) {
		fmt.Println("⚡ Significant change detected!")

		oldContent := record.PreviousFileContent

		// Compute delta between old and new versions
		delta := diffalgo.ComputeDelta(record.File, record.File, oldContent, newContent)

		// Print delta in readable format
		fmt.Println("===== Delta =====")
		for _, d := range delta {
			fmt.Printf("File: %s | Line: %d | Type: %s | Content: %v\n",
				d.FilePath, d.LineNumber, d.Type, d.Content)
			for _, c := range d.CharDiff {
				fmt.Printf("    CharDiff: %s -> %q\n", c.Type, c.Text)
			}
		}

		// Convert delta to JSON string
		deltaBytes, err := json.MarshalIndent(delta, "", "  ")
		if err != nil {
			fmt.Println("❌ Error serializing delta:", err)
			return
		}

		// Save diff as blob (timestamped blob file)
		deltaJSON := string(deltaBytes)
		blobPath := w.Blob.CreateBlobFromContent(deltaJSON)
		fmt.Println("✅ Delta blob created at:", blobPath)

		// save to history
		data := types.FileRecord{
			File: path,
			Type: "delta",
			Action: "write",
			Blob: blobPath,
			IsBlobType: true,
			Timestamp: time.Now(),
		}

		err = w.History.Create(path, data)
		if err != nil{
			log.Fatal("writeEvent:",err)
		}
		fmt.Println("history created for write delta!")

		// Update file record
		record.PrevSize = record.CurrentSize
		record.PreviousFileContent = newContent
		record.Timestamp = time.Now()

		// Save updated state
		w.State[path] = record
		fmt.Println("💾 File record updated")

	} else {
		fmt.Println("ℹ️ No significant change detected")
	}
}
