package updatecmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Update struct {
	UpdateCmd *cobra.Command
}

func NewUpdateCmd(use, short string) *Update {
	return &Update{
		UpdateCmd: &cobra.Command{
			Use:   use,
			Short: short,
		},
	}
}

func (u *Update) Run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Error: Please provide a new URL.\nUsage: rec update <new-url>")
		return
	}

	newURL := args[0]
	configPath := filepath.Join(".rec", "config.json")

	// Check if .rec/config.json exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("❌ No repository found. Run 'rec set <url>' first.")
		return
	}

	// Load existing config
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

	// Update repository.remote
	repo, ok := config["repository"].(map[string]interface{})
	if !ok {
		fmt.Println("❌ Invalid config format. Run 'rec set <url>' to reset.")
		return
	}

	oldURL, _ := repo["remote"].(string)
	repo["remote"] = newURL
	config["repository"] = repo

	// Save updated config
	updatedData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error encoding config:", err)
		return
	}
	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		fmt.Println("Error writing config.json:", err)
		return
	}

	fmt.Println("✅ Remote URL updated successfully!")
	fmt.Printf("→ Old: %s\n→ New: %s\n", oldURL, newURL)
}
