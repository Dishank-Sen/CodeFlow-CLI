package cli

import (
	"exp1/utils"
	"fmt"

	"github.com/spf13/cobra"
)

func Root() *cobra.Command{
	var rootCmd = &cobra.Command{
		Use: "rec",
		Short: "rec is a simple version control system",
		Long: "rec is a version control system built in Go. It captures code changes.",
		PersistentPreRunE: persistentPreRunE,
	}

	// loop which register all the commands
	for _, cmd := range Registered{
		c := cmd()
		rootCmd.AddCommand(c)
	}

	return rootCmd
}

func persistentPreRunE(cmd *cobra.Command, args []string) error{
	if cmd.Name() == "init"{
		return nil
	}
	
	// if .rec is not created prompt user to run init command
	if !utils.CheckDirExist(".rec"){
		err := "not a rec repository, run 'rec init' to initialize a empty rec repository"
		return fmt.Errorf(err)
	}
	
	return nil
}