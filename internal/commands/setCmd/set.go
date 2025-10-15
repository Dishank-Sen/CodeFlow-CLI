package setcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Set struct {
	SetCmd *cobra.Command
}

func NewSetCmd(use, short string) *Set {
	s := &Set{
		SetCmd: &cobra.Command{
			Use:   use,
			Short: short,
			Run:   nil, // we'll set below
		},
	}

	// Define flags
	s.SetCmd.Flags().String("username", "", "Git username")
	s.SetCmd.Flags().String("reponame", "", "Repository name")
	s.SetCmd.Flags().String("remoteUrl", "", "Remote repository URL")

	// Assign run function
	s.SetCmd.Run = s.Run

	return s
}

func (s *Set) Run(cmd *cobra.Command, args []string) {
	// Read flag values
	userName, _ := cmd.Flags().GetString("username")
	repoName, _ := cmd.Flags().GetString("reponame")
	remoteUrl, _ := cmd.Flags().GetString("remoteUrl")

	configPath := filepath.Join(".rec", "config.json")

	// Ensure .rec directory exists
	if _, err := os.Stat(".rec"); os.IsNotExist(err) {
		if err := os.Mkdir(".rec", 0755); err != nil {
			fmt.Println("❌ Failed to create .rec directory:", err)
			return
		}
	}

	// Load existing config if present
	config := make(map[string]interface{})
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err == nil {
			json.Unmarshal(data, &config)
		}
	}

	// Update only provided fields
	if userName != "" {
		config["userName"] = userName
	}
	if repoName != "" {
		config["repoName"] = repoName
	}
	if remoteUrl != "" {
		config["remoteUrl"] = remoteUrl
	}

	// Write back to config.json
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("❌ Error encoding config:", err)
		return
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Println("❌ Error writing config.json:", err)
		return
	}

	fmt.Println("✅ Repository configuration updated successfully!")
	if userName != "" {
		fmt.Println("→ userName:", userName)
	}
	if repoName != "" {
		fmt.Println("→ repoName:", repoName)
	}
	if remoteUrl != "" {
		fmt.Println("→ remoteUrl:", remoteUrl)
	}
}
