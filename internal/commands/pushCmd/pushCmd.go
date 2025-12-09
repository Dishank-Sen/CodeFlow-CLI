package pushcmd

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
)

func NewPushCmd(use, short string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Run: Run,
	}
}

func Run(cmd *cobra.Command, args []string) {
	configPath := filepath.Join(".rec", "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("No .rec/config.json found. Run 'rec init' and 'rec set --remoteUrl <url>' first.")
		return
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading config.json:", err)
		return
	}

	config := make(map[string]interface{})
	if err := json.Unmarshal(data, &config); err != nil {
		fmt.Println("Error parsing config.json:", err)
		return
	}

	remote, ok := config["remoteUrl"].(string) // use root remoteUrl
	if !ok || remote == "" {
		fmt.Println("No remote URL found. Run 'rec set --remoteUrl <url>' to set it.")
		return
	}

	res, err := PushTriggered(remote)
	if err != nil {
		log.Fatal("error occurs (pushCmd):", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("error reading response:", err)
	}

	fmt.Println("response status:", res.Status)
	fmt.Println("response body:", string(body))
}

func PushTriggered(remote string) (*http.Response, error){
	// get pipe reader and writer
	pr, pw := io.Pipe()

	// get a writer to the pipe
	zipWriter := zip.NewWriter(pw)

	go func(){
		filepath.Walk(".rec/history", func(path string, info fs.FileInfo, err error) error {
			if info.IsDir(){
				return nil
			}

			f, err := os.Open(path)
			if err != nil{
				return err
			}

			rel, err := filepath.Rel(".rec", path)
			if err != nil{
				return err
			}

			w, err := zipWriter.Create(rel)
			
			io.Copy(w, f)
			return nil
		})
		zipWriter.Close()
		pw.Close()
	}()

	return http.Post(remote, "application/zip", pr)
}