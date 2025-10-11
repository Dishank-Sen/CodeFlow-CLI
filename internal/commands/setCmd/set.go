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
	return &Set{
		SetCmd: &cobra.Command{
			Use:   use,
			Short: short,
		},
	}
}

func (s *Set) Run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Please provide a URL.\nUsage: rec set <url>")
		return
	}

	url := args[0]
	configPath := filepath.Join(".rec", "config.json")

	// Ensure .rec directory exists
	if _, err := os.Stat(".rec"); os.IsNotExist(err) {
		if err := os.Mkdir(".rec", 0755); err != nil {
			fmt.Println("Failed to create .rec directory:", err)
			return
		}
	}

	// Load existing config if present
	config := make(map[string]interface{})
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Println("Error reading config.json:", err)
			return
		}
		json.Unmarshal(data, &config)
	}

	// Update repository section
	repo, ok := config["repository"].(map[string]interface{})
	if !ok {
		repo = make(map[string]interface{})
	}
	repo["remote"] = url
	config["repository"] = repo

	// Write back to file
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error encoding config:", err)
		return
	}
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		fmt.Println("Error writing config.json:", err)
		return
	}

	fmt.Println("✅ Remote URL set successfully!")
	fmt.Println("→", url)
}
