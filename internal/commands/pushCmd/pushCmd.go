package pushcmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

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

	userName, _ := config["userName"].(string)
	repoName, _ := config["repoName"].(string)

	fmt.Printf("🚀 Pushing repository '%s' by user '%s' to remote: %s\n", repoName, userName, remote)

	if err := PushTriggered(remote); err != nil {
		fmt.Println("❌ Push failed:", err)
		return
	}

	fmt.Println("✅ Push completed successfully!")
}


// PushTriggered zips .rec, and uploads to Express server
func PushTriggered(remote string) error {
	fmt.Println("📦 Collecting project files...")

	// Zip folder to disk
	zipPath := "project.zip"
	buf, err := zipFolder(".rec")
	if err != nil {
		return fmt.Errorf("failed to zip folder: %w", err)
	}

	if err := os.WriteFile(zipPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write zip to disk: %w", err)
	}

	// Prepare multipart request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", `form-data; name="file"; filename="project.zip"`)
	h.Set("Content-Type", "application/zip")
	part, err := writer.CreatePart(h)
	if err != nil {
		os.Remove(zipPath)
		return err
	}

	if _, err := io.Copy(part, buf); err != nil {
		os.Remove(zipPath)
		return err
	}
	writer.Close()

	// Send request
	req, err := http.NewRequest("POST", remote, body)
	if err != nil {
		os.Remove(zipPath)
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		os.Remove(zipPath)
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Server response:", string(respBody))

	// Check if response is suitable (customize this condition)
	if resp.StatusCode != 200 || !bytes.Contains(respBody, []byte("success")) {
		fmt.Println("⚠️ Response not suitable, removing zip file.")
		os.Remove(zipPath)
		return fmt.Errorf("server response unsuitable")
	}

	// Otherwise, remove zip anyway if you don't need it
	os.Remove(zipPath)

	return nil
}

// CollectProjectFiles copies project files into .rec/files
func CollectProjectFiles() error {
	projectRoot, _ := os.Getwd()
	recDir := filepath.Join(projectRoot, ".rec")
	filesDir := filepath.Join(recDir, "files")

	os.MkdirAll(filesDir, os.ModePerm)

	ignoreFile := filepath.Join(projectRoot, ".recignore")
	ignorePatterns := []string{}
	if data, err := os.ReadFile(ignoreFile); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				ignorePatterns = append(ignorePatterns, line)
			}
		}
	}

	return filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == recDir {
			return filepath.SkipDir
		}

		rel, _ := filepath.Rel(projectRoot, path)
		for _, pattern := range ignorePatterns {
			match, _ := filepath.Match(pattern, rel)
			if match || strings.HasPrefix(rel, pattern) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		dest := filepath.Join(filesDir, rel)
		os.MkdirAll(filepath.Dir(dest), os.ModePerm)
		copyFile(path, dest)
		fmt.Println("Added:", rel)
		return nil
	})
}

func copyFile(src, dst string) error {
	in, _ := os.Open(src)
	defer in.Close()
	out, _ := os.Create(dst)
	defer out.Close()
	_, err := io.Copy(out, in)
	return err
}

func zipFolder(folderPath string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		// close immediately (not deferred)
		defer file.Close()

		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func RemoveFile(path string){
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}