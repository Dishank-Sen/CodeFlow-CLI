package recorder

import "time"

type FileState struct {
	CurrentSize int64  `json:"currentSize"`
	PrevSize    int64  `json:"prevSize"`
	Hash        string `json:"hash"`
}

type SnapshotFile struct {
	File      string    `json:"file"`
	Type      string    `json:"type"`
	Blob      string    `json:"blob"`
	Timestamp time.Time `json:"timestamp"`
}

type RemoveFile struct {
	File      string    `json:"file"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

type Change struct {
	Action string `json:"action"`
	Name   string `json:"name"`
}

type DeltaFile struct {
	File      string    `json:"file"`
	Type      string    `json:"type"`
	Base      string    `json:"base"`
	Parent    string    `json:"parent"`
	Blob      string    `json:"blob"`
	Timestamp time.Time `json:"timestamp"`
	Changes   []Change  `json:"changes"`
}
