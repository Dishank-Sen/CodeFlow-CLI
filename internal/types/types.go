package types

import "time"

type FileRecord struct {
	File      string    `json:"file"`  // file path
	Type      string    `json:"type,omitempty"`       // e.g. "snapshot", "delta"
	Action    string    `json:"action"`     // optional: e.g. "create", "write", "delete", "remove",
	Content   string    `json:"content,omitempty"`
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

type Repository struct{
	UserName string `json:"username"`
	RemoteUrl string `json:"remoteUrl"`
}

type Recorder struct{
	DebounceTime int64
}

type Config struct {
	Repository Repository
	Recorder Recorder
}

type Node struct {
    Name     string     `json:"name"`
    Path     string     `json:"path"`               // absolute or repo-relative
    IsDir    bool       `json:"isDir"`
    Size     int64      `json:"size,omitempty"`     // bytes; 0 for dirs
	CreateTime time.Time `json:"createTime,omitempty"`
    ModTime  time.Time  `json:"modTime,omitempty"`
    Children []*Node    `json:"children,omitempty"` // nil when no children -> omitted in JSON
}

type FileTree struct{
	Files []*Node `json:"files"`
}