package main

import (
	"context"
	"exp1/cli"
	"exp1/utils/log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)


func main(){
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Error(ctx, stop, err.Error())	
	}

	rootCmd := cli.Root(ctx)
	rootCmd.Execute()
}