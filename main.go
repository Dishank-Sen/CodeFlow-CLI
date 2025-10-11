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
	initCmd.InitCmd.Run = initCmd.Run
	Register(initCmd.InitCmd, rootCmd)
}

func RegisterSetCmd(rootCmd *cmd.Root) {
	use := "set"
	short := "Set repository URL"
	setCommand := setcmd.NewSetCmd(use, short)
	setCommand.SetCmd.Run = setCommand.Run
	Register(setCommand.SetCmd, rootCmd)
}

func RegisterUpdateCmd(rootCmd *cmd.Root) {
	use := "update"
	short := "Update existing repository URL"
	updateCommand := updatecmd.NewUpdateCmd(use, short)
	updateCommand.UpdateCmd.Run = updateCommand.Run
	Register(updateCommand.UpdateCmd, rootCmd)
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
	startCmd.StartCmd.Run = startCmd.Run
	Register(startCmd.StartCmd, rootCmd)
}

func main(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	use := "rec"
	short := "rec is a simple version control system"
	long := "rec is a toy version control system built in Go. It mimics some functionality of git for learning purposes."

	rootCmd := cmd.NewRootCmd(use, short, long)
	RegisterAllCmd(rootCmd)
	rootCmd.Execute()
}