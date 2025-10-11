package delta

import (
	"encoding/json"
	"exp1/internal/types"
	"fmt"
	"os"
	"path/filepath"
)

type Change struct {
	Action string `json:"action"`
	Name   string `json:"name"`
}

type Delta struct{

}

func NewDelta() *Delta{
	return &Delta{}
}

func (d *Delta) Create(path string, data any) error{
	// Create history directory for the file
	fileDir, err := d.getDirForFile(path)
	if err != nil {
		return err
	}

	// Generate the history file name
	title := d.getFileTitle(data)
	filePath := filepath.Join(fileDir, title)

	// Marshal data to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal nonblob history file: %w", err)
	}

	// Save JSON file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write nonblob history file: %w", err)
	}

	return nil
}

func (d *Delta) getDirForFile(path string) (string, error) {
	historyDir := filepath.Join(".rec", "history", path)
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create history dir: %w", err)
	}
	return historyDir, nil
}

func (d *Delta) getFileTitle(data any) string {
	if fileData, ok := data.(types.FileRecord); ok {
		return fileData.Timestamp.Format("20060102_150405") + ".json"
	}
	return "unknown.json"
}