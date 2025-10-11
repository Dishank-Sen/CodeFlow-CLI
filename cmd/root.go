package cmd

import (
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

func (r *Root) init(){

}
