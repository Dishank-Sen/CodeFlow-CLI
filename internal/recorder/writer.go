package recorder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Writer struct {
    root string
}

func NewWriter(root string) *Writer {
    return &Writer{root: root}
}

func (w *Writer) SaveBlob(path string) (string, error) {
    // Compute hash
    fmt.Println("computing hash..")
    hash, err := HashFile(path)
    if err != nil {
        fmt.Println("writer.go :", err)
        return "", err
    }

    // If file is empty, return early
    if hash == "" {
        fmt.Println("file is empty")
        return "", nil
    }

    // Construct blob path
    blobPath := filepath.Join(".rec", "blob", hash[:2], hash)
    if err := os.MkdirAll(filepath.Dir(blobPath), 0755); err != nil {
        return "", err
    }

    // Read file content
    content, err := os.ReadFile(path)
    if err != nil {
        fmt.Println("writer.go content error:", err)
        return "", err
    }

    // Write content to blob
    if err := os.WriteFile(blobPath, content, 0644); err != nil {
        return "", err
    }

    return blobPath, nil
}

func (w *Writer) SaveSnapshotHistory(path string, snapshotFile SnapshotFile) error{
    // Create history directory for the file
    historyDir := filepath.Join(".rec", "history", path)
    if err := os.MkdirAll(historyDir, 0755); err != nil {
        return fmt.Errorf("failed to create history dir: %w", err)
    }

    // Title can be timestamp-based (unique)
    title := snapshotFile.Timestamp.Format("20060102_150405") + ".json"
    historyPath := filepath.Join(historyDir, title)

    // Marshal snapshotFile to JSON
    data, err := json.MarshalIndent(snapshotFile, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal snapshotFile: %w", err)
    }

    // Save JSON to file
    if err := os.WriteFile(historyPath, data, 0644); err != nil {
        return fmt.Errorf("failed to write history file: %w", err)
    }

    fmt.Println("Snapshot history saved at:", historyPath)
    return nil
}

func (w *Writer) SaveFileRemoveHistory(path string, removeFile RemoveFile) error{
    // Create history directory for the file
    historyDir := filepath.Join(".rec", "history", path)
    if err := os.MkdirAll(historyDir, 0755); err != nil {
        return fmt.Errorf("failed to create history dir: %w", err)
    }

    // Title can be timestamp-based (unique)
    title := removeFile.Timestamp.Format("20060102_150405") + ".json"
    historyPath := filepath.Join(historyDir, title)

    // Marshal snapshotFile to JSON
    data, err := json.MarshalIndent(removeFile, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal removeFile: %w", err)
    }

    // Save JSON to file
    if err := os.WriteFile(historyPath, data, 0644); err != nil {
        return fmt.Errorf("failed to write history file: %w", err)
    }

    fmt.Println("Snapshot history for file removal saved at:", historyPath)
    return nil
}