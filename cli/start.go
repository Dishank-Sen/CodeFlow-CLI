package cli

import (
	"exp1/internal/commands/startCmd/events"
	"exp1/internal/commands/startCmd/watcher"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func init(){
	Register("start", Start)
}

func Start() *cobra.Command{
	return &cobra.Command{
		Use: "start",
		Short: "starts recording file changes",
		RunE: startRunE,
	}
}

func startRunE(cmd *cobra.Command, args []string) error{
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
	return nil
}