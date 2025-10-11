package recorder

import (
	"exp1/internal/recorder/blob"
	"exp1/internal/recorder/history"
)

type SnapshotManager interface {
    CreateSnapshot(path string) error
    RestoreSnapshot(path string, version string) error
    Diff(oldPath, newPath string) (*Delta, error)
}

// type Snapshot struct {
//     ID      string
//     Path    string
//     Version string
//     Hash    string
//     Size    int64
// }

type Snapshot struct{
    blob *blob.Blob
    history *history.History
    // index
    // config
}

func NewSnapshot() *Snapshot{
    blob := blob.NewBlob()
    history := history.NewHistory()
    return &Snapshot{
        blob: blob,
        history: history,
    }
}