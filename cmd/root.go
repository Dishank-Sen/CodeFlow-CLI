package cmd

import (
	"exp1/internal/commands/initCmd"
	"os"

	"github.com/spf13/cobra"
)

type Root struct{
	RootCmd *cobra.Command
}

func NewRootCmd(use string, short string, long string) *Root{
	var rootCmd = &cobra.Command{
		Use: use,
		Short: short,
		Long: long,
		PersistentPreRunE: ensureRecExist,
	}
	return &Root{
		RootCmd: rootCmd,
	}
}

func (r *Root) Register(cmd *cobra.Command){
	r.RootCmd.AddCommand(cmd)
}

func (r *Root) Execute() {
	err := r.RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ensureRecExist(cmd *cobra.Command, args []string) error{
	if cmd.Name() == "init"{
		return nil
	}
	
	// run the init command logic
	return initCmd.Run(cmd, args)
}