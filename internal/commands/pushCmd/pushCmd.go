package pushcmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type Push struct {
	PushCmd *cobra.Command
}

func NewPushCmd(use, short string) *Push {
	return &Push{
		PushCmd: &cobra.Command{
			Use:   use,
			Short: short,
		},
	}
}

func (p *Push) Run(cmd *cobra.Command, args []string) {
	configPath := filepath.Join(".rec", "config.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("❌ No .rec/config.json found. Run 'rec init' and 'rec set <url>' first.")
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

	repo, ok := config["repository"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Invalid config format.")
		return
	}

	remote, ok := repo["remote"].(string)
	if !ok || remote == "" {
		fmt.Println("❌ No remote URL found. Run 'rec set <url>' to set it.")
		return
	}

	fmt.Println("🚀 Pushing repository to remote:", remote)

	if err := p.PushTriggered(remote); err != nil {
		fmt.Println("❌ Push failed:", err)
		return
	}

	fmt.Println("✅ Push completed successfully!")
}

// PushTriggered zips .rec, and uploads to Express server
func (p *Push) PushTriggered(remote string) error {
	fmt.Println("📦 Collecting project files...")

	if err := p.CollectProjectFiles(); err != nil {
		return fmt.Errorf("error collecting project files: %w", err)
	}
	fmt.Println("✅ All project files copied to .rec/files successfully.")

	// Zip folder into memory
	zipBuf, err := zipFolder(".rec")
	if err != nil {
		return fmt.Errorf("failed to zip folder: %w", err)
	}

	// Prepare multipart request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := textproto.MIMEHeader{}
	h.Set("Content-Disposition", `form-data; name="file"; filename="project.zip"`)
	h.Set("Content-Type", "application/zip")
	part, err := writer.CreatePart(h)
	if err != nil {
		return err
	}

	// Copy zip from memory directly to multipart
	if _, err := io.Copy(part, zipBuf); err != nil {
		return err
	}

	writer.Close() // finalize multipart

	// Send request
	req, err := http.NewRequest("POST", remote, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Server response:", string(respBody))

	return nil
}

// CollectProjectFiles copies project files into .rec/files
func (p *Push) CollectProjectFiles() error {
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

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(folderPath, path)
		file, _ := os.Open(path)
		defer file.Close()
		w, _ := zipWriter.Create(rel)
		_, err = io.Copy(w, file)
		return err
	})

	if err != nil {
		return nil, err
	}

	zipWriter.Close()
	return buf, nil
}
