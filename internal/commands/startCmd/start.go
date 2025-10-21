package startCmd

import (
	"exp1/internal/commands/startCmd/events"
	"exp1/internal/commands/startCmd/watcher"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type Start struct{
	StartCmd *cobra.Command
}

func NewStartCmd(use string, short string) *cobra.Command{
	return &cobra.Command{
		Use: use,
		Short: short,
		Run: Run,
	}
}

func Run(cmd *cobra.Command, args []string){
   	w := watcher.NewWatcher()
    ev := events.NewEvents(w)
    w.SetEvents(ev)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\n🛑 Terminating... flushing unsaved deltas.")
		ev.Flush() // flush unsaved changes
		os.Exit(0)
	}()

    w.Start()
}