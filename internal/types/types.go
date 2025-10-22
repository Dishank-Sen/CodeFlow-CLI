package types

import "time"


type Change struct {
	Action string `json:"action"`
	Name   string `json:"name"`
}

type FileRecord struct {
	File      string    `json:"file"`  // file path
	Type      string    `json:"type,omitempty"`       // e.g. "snapshot", "delta"
	Action    string    `json:"action"`     // optional: e.g. "create", "write", "delete", "remove",
	Blob      string    `json:"blob,omitempty"`
	IsBlobType bool     `json:"isBlobType,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	OldPath   string    `json:"oldPath,omitempty"`
	NewPath   string    `json:"newPath,omitempty"`

	// Optional file state
	CurrentSize int64  `json:"currentSize,omitempty"`
	PrevSize    int64  `json:"prevSize,omitempty"`
	PreviousFileContent        string `json:"previousFileContent,omitempty"`
}

// CharDiff represents character-level changes within a line
type CharDiff struct {
	Type string `json:"Type"` // "Equal", "Insert", "Delete"
	Text string `json:"Text"`
}

// LineChange represents line-level changes
type LineChange struct {
	FilePath   string     `json:"FilePath"`   // which file this change belongs to
	LineNumber int        `json:"LineNumber"` // index in the old file
	Type       string     `json:"Type"`       // "add", "delete", "replace"
	Content    []string   `json:"Content"`    // content for added/replaced lines
	CharDiff   []CharDiff `json:"CharDiff"`   // optional: intra-line diff for replace
}

type Config struct {
	Repository struct {
		Remote string `json:"remote"`
	} `json:"repository"`
}