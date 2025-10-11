package startCmd

import (
	"exp1/internal/commands/startCmd/events"
	"exp1/internal/commands/startCmd/watcher"

	"github.com/spf13/cobra"
)

type Start struct{
	StartCmd *cobra.Command
}

func NewStartCmd(use string, short string) *Start{
	var startCmd = &cobra.Command{
		Use: use,
		Short: short,
	}
	return &Start{
		StartCmd: startCmd,
	}
}

func (s *Start) Run(cmd *cobra.Command, args []string){
   	w := watcher.NewWatcher()
    ev := events.NewEvents(w)
    w.SetEvents(ev)
    w.Start()
}