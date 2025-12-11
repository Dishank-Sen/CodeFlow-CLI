package rename

import (
	"context"
	roottimeline "exp1/internal/recorder/root-timeline"
	"exp1/internal/types"
	"exp1/pkg/interfaces"
	"os"
	"path/filepath"
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
	newName := filepath.Base(r.NewPath)
	oldName := filepath.Base(r.OldPath)
	info, err := os.Stat(r.NewPath)
	if err != nil{
		return err
	}

	var data = types.Rename{
		NewPath: r.NewPath,
		NewName: newName,
		Action: "rename",
		OldPath: r.OldPath,
		OldName: oldName,
		IsDir: info.IsDir(),
		Size: info.Size(),
		RenameTime: time.Now(),
	}

	// add file to .rec/history
	return roottimeline.Save(data)
}
