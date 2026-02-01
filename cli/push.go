package cli

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"exp1/cli/utils"
	"exp1/internal/types"
	"exp1/utils/log"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	Register("push", Push)
}

func Push() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "pushes all the snapshot and deltas to server",
		RunE:  pushRunE,
	}
}

func pushRunE(cmd *cobra.Command, args []string) error {
	configPath := filepath.Join(".codeflow", "config.json")
	parentCtx := cmd.Context()
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Info(parentCtx, "no config file exist")
		log.Info(parentCtx, "creating default config file.")

		// create a default config file
		err := utils.CreateConfig(ctx, cancel, false)
		if err != nil {
			return err
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var config types.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	remoteUrl := config.Repository.RemoteUrl

	if strings.TrimSpace(remoteUrl) == "" {
		return fmt.Errorf("no remote url found, run codeflowset -r <remoteUrl> to set it.")
	}

	userName, repoName, err := parseRemoteURL(remoteUrl)
	fmt.Printf("username: %s and repoName: %s", userName, repoName)
	if err != nil {
		return err
	}

	endpointUrl := "http://localhost:3000/api/v1/push"

	res, err := Trigger(userName, repoName, endpointUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fmt.Println("response status:", res.Status)
	fmt.Println("response body:", string(body))
	statusMsg := fmt.Sprintf("response status: %s", res.Status)
	log.Info(parentCtx, statusMsg)

	bodyMsg := fmt.Sprintf("response body: %s", string(body))
	log.Info(parentCtx, bodyMsg)

	return nil
}

func parseRemoteURL(remoteUrl string) (username string, repoName string, err error) {
	// Ensure URL has a scheme for proper parsing
	if !strings.Contains(remoteUrl, "://") {
		remoteUrl = "https://" + remoteUrl
	}

	u, err := url.Parse(remoteUrl)
	if err != nil {
		return "", "", err
	}

	// Extract path segments
	segments := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(segments) != 2 {
		return "", "", errors.New("invalid URL format: expected /<username>/<repo>.codeflow")
	}

	username = segments[0]
	repo := segments[1]

	if !strings.HasSuffix(repo, ".codeflow") {
		return "", "", errors.New("invalid repo name: missing .codeflowsuffix")
	}

	repoName = strings.TrimSuffix(repo, ".codeflow")
	if repoName == "" {
		return "", "", errors.New("empty repo name")
	}

	return username, repoName, nil
}

func Trigger(userName, repoName, endpointUrl string) (*http.Response, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()
		defer writer.Close()

		// --- metadata part ---
		metaPart, _ := writer.CreateFormField("metadata")
		metadata := types.Metadata{
			UserName: userName,
			RepoName: repoName,
		}
		metadataBytes, _ := json.Marshal(metadata)
		metaPart.Write(metadataBytes)

		zipFiles(writer, "history", "history.zip", ".codeflow/history")
		zipFiles(writer, "fileTree", "fileTree.zip", ".codeflow/files")
		zipFiles(writer, "root-timeline", "root-timeline.zip", ".codeflow/root-timeline")
	}()

	req, _ := http.NewRequest("POST", endpointUrl, pr)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	return client.Do(req)
}

func zipFiles(writer *multipart.Writer, fieldname string, filename string, dirPath string) {
	zipPart, _ := writer.CreateFormFile(fieldname, filename)
	zipWriter := zip.NewWriter(zipPart)

	filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		f, _ := os.Open(path)
		defer f.Close()

		rel, _ := filepath.Rel(".codeflow", path)
		w, _ := zipWriter.Create(rel)
		io.Copy(w, f)
		return nil
	})

	zipWriter.Close()
}
