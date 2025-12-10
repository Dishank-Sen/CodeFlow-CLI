package cli

import (
	"context"
	initfiles "exp1/cli/initFiles"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init(){
	Register("init", Init)
}

func Init() *cobra.Command{
	return &cobra.Command{
		Use: "init",
		Short: "nitialize a new rec repository",
		RunE: initRunE,
	}
}

func initRunE(cmd *cobra.Command, args []string) error{
	ctx, cancel := context.WithCancel(cmd.Context())

    rootFolder := ".rec"

    folders := getFolders(rootFolder)

    if  folderExist(rootFolder){
		absPath, _ := filepath.Abs(rootFolder)
        fmt.Printf("Reinitialized existing rec repo in %s\n", absPath)
		return nil
	}

    // Create directories
    err := createDir(folders)
	if err != nil{
		log.Fatal(err)
		return err
	}

    // Create files
    err = createFiles(ctx, cancel)
	if err != nil{
		return err
	}

    absPath, _ := filepath.Abs(rootFolder)
    fmt.Printf("Initialized empty rec repo in %s\n", absPath)
	return nil
}

func getFolders(rootFolder string) []string{
	return []string{
        filepath.Join(rootFolder, "blob"),
        filepath.Join(rootFolder, "history"),
        filepath.Join(rootFolder, "index"),
        filepath.Join(rootFolder, "files"),
		filepath.Join(rootFolder, "root-timeline"),
    }
}

func folderExist(rootFolder string) bool{
	_, err := os.Stat(rootFolder)
	return err == nil
}

func createDir(folders []string) error{
	for _, folder := range folders {
        if err := os.MkdirAll(folder, 0755); err != nil {
            return fmt.Errorf("error creating folders: %v", err)
        }
    }
	return nil
}

func createFiles(ctx context.Context, cancel context.CancelFunc) error{
	for _, f := range initfiles.InitFiles{
		err := f(ctx, cancel)
		if err != nil{
			return err
		}
	}
	return nil
}