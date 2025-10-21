package main

import (
	"exp1/cmd"
	"exp1/internal/commands/initCmd"
	pushcmd "exp1/internal/commands/pushCmd"
	setcmd "exp1/internal/commands/setCmd"
	"exp1/internal/commands/startCmd"
	updatecmd "exp1/internal/commands/updateCmd"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func Register(cmd *cobra.Command, rootCmd *cmd.Root){
	rootCmd.Register(cmd)
}

func RegisterAllCmd(rootCmd *cmd.Root){
	RegisterInitCmd(rootCmd)
	RegisterSetCmd(rootCmd)
	RegisterStartCmd(rootCmd)
	RegisterUpdateCmd(rootCmd)
	RegisterPushCmd(rootCmd)
}

func RegisterInitCmd(rootCmd *cmd.Root){
	use := "init"
	short := "Initialize a new rec repository"
	initCmd := initCmd.NewInitCmd(use, short)
	Register(initCmd, rootCmd)
}

func RegisterSetCmd(rootCmd *cmd.Root) {
	use := "set"
	short := "Set repository URL"
	setCommand := setcmd.NewSetCmd(use, short)
	Register(setCommand, rootCmd)
}

func RegisterUpdateCmd(rootCmd *cmd.Root) {
	use := "update"
	short := "Update existing repository URL"
	updateCommand := updatecmd.NewUpdateCmd(use, short)
	Register(updateCommand, rootCmd)
}

func RegisterPushCmd(rootCmd *cmd.Root) {
	use := "push"
	short := "Push recorded changes to the remote repository"
	pushCommand := pushcmd.NewPushCmd(use, short)
	pushCommand.PushCmd.Run = pushCommand.Run
	Register(pushCommand.PushCmd, rootCmd)
}

func RegisterStartCmd(rootCmd *cmd.Root){
	use := "start"
	short := "Starts recording file changes"
	startCmd := startCmd.NewStartCmd(use, short)
	Register(startCmd, rootCmd)
}

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	use := "rec"
	short := "rec is a simple version control system"
	long := "rec is a version control system built in Go. It captures code changes."

	rootCmd := cmd.NewRootCmd(use, short, long)
	RegisterAllCmd(rootCmd)
	rootCmd.Execute()
}