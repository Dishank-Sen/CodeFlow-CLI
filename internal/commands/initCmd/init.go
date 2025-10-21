package initCmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Init struct{
	InitCmd *cobra.Command
}

func NewInitCmd(use string, short string) *Init{
	var initCmd = &cobra.Command{
		Use: use,
		Short: short,
	}
	return &Init{
		InitCmd: initCmd,
	}
}

func (i *Init) Run(cmd *cobra.Command, args []string){
    rootFolder := ".rec"

    folders := i.getFolders(rootFolder)

    files := i.getFiles(rootFolder)

    if  i.folderExist(rootFolder){
		absPath, _ := filepath.Abs(rootFolder)
        fmt.Printf("Reinitialized existing rec repo in %s\n", absPath)
		return
	}

    // Create directories
    err := i.createDir(folders)
	if err != nil{
		log.Fatal(err)
		return
	}

    // Create files
    err = i.createFiles(files)
	if err != nil{
		log.Fatal(err)
		return
	}

    absPath, _ := filepath.Abs(rootFolder)
    fmt.Printf("Initialized empty rec repo in %s\n", absPath)
}

func (i *Init) getFolders(rootFolder string) []string{
	return []string{
        filepath.Join(rootFolder, "blob"),
        filepath.Join(rootFolder, "history"),
        filepath.Join(rootFolder, "index"),
        filepath.Join(rootFolder, "files"),
		filepath.Join(rootFolder, "root-timeline"),
    }
}

func (i *Init) getFiles(rootFolder string) []string{
	return []string{
        filepath.Join(rootFolder, "index", "index.json"),
        filepath.Join(rootFolder, "config.json"),
    }
}

func (i *Init) folderExist(rootFolder string) bool{
	if _, err := os.Stat(rootFolder); err == nil {
        return true
    }
	return false
}

func (i *Init) createDir(folders []string) error{
	for _, folder := range folders {
        if err := os.MkdirAll(folder, 0755); err != nil {
            return fmt.Errorf("error creating folders: %v", err)
        }
    }
	return nil
}

func (i *Init) createFiles(files []string) error{
	for _, file := range files {
        f, err := os.Create(file)
        if err != nil {
            return fmt.Errorf("error creating file: %v", err)
        }
        f.Close()
    }
	return nil
}