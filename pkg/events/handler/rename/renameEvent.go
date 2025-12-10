package rename

import (
	"context"
	"exp1/pkg/interfaces"
	roottimeline "exp1/internal/recorder/root-timeline"
	"exp1/internal/types"
	"time"
)

type Rename struct{
	OldPath string
	NewPath string
	Watcher interfaces.IWatcher
	Ctx context.Context
}

func NewRename(ctx context.Context, oldPath string, newpath string, watcher interfaces.IWatcher) *Rename{
	return &Rename{
		OldPath: oldPath,
		NewPath: newpath,
		Watcher: watcher,
		Ctx: ctx,
	}
}

func (r *Rename) Trigger() error{

	var data = types.FileRecord{
		File: r.NewPath,
		Action: "rename",
		NewPath: r.NewPath,
		OldPath: r.OldPath,
		Timestamp: time.Now(),
	}

	// add file to .rec/history
	return roottimeline.Save(data)
}
