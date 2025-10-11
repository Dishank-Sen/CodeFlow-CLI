package recorder

import (
	"fmt"
	"time"
)

type Manager struct {
    Root   string
    Writer *Writer
    Reader *Reader
}

func NewManager(root string) *Manager {
    return &Manager{
        Root:   root,
        Writer: NewWriter(root),
        Reader: NewReader(root),
    }
}

func (m *Manager) CreateSnapshot(path string) error {
    // first save the blob
    blobPath, err := m.Writer.SaveBlob(path)
    if err != nil{
        fmt.Println(err)
        return err
    }

    // if blobpath is empty return safely
    if blobPath == ""{
        return nil
    }

    // make snapshot object
    snapshotFile := SnapshotFile{
        File: path,
        Type: "snapshot",
        Blob: blobPath,
        Timestamp: time.Now(),
    }
    
    // save snapshot history
    err = m.Writer.SaveSnapshotHistory(path, snapshotFile)
    if err != nil{
        fmt.Println(err)
        return err
    }

    return nil
}

func (m *Manager) FileRemoveSnapshot(path string) error{
    removeFile := RemoveFile{
        File: path,
        Type: "remove",
        Timestamp: time.Now(),
    }

    // save file remove snapshot
    err := m.Writer.SaveFileRemoveHistory(path, removeFile)
    if err != nil{
        fmt.Println(err)
        return err
    }

    return nil
}

func (m *Manager) RestoreSnapshot(path string, version string) ([]byte, error) {
    return m.Reader.Load(path, version)
}
