package initCmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
)

func NewInitCmd(use string, short string) *cobra.Command{
	return &cobra.Command{
		Use: use,
		Short: short,
		RunE: Run,
	}
}

func Run(cmd *cobra.Command, args []string) error{
    rootFolder := ".rec"

    folders := getFolders(rootFolder)

    files := getFiles(rootFolder)

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
    err = createFiles(files)
	if err != nil{
		log.Fatal(err)
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

func getFiles(rootFolder string) []string{
	return []string{
        filepath.Join(rootFolder, "index", "index.json"),
        filepath.Join(rootFolder, "config.json"),
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

func createFiles(files []string) error{
	for _, file := range files {
        f, err := os.Create(file)
        if err != nil {
            return fmt.Errorf("error creating file: %v", err)
        }
        f.Close()
    }
	return nil
}