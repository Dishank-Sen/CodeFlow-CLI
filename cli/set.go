package cli

import (
	"encoding/json"
	"exp1/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init(){
	Register("set", Set)
}

func Set() *cobra.Command{
	SetCmd := &cobra.Command{
		Use:   "set",
		Short: "Set repository specific info",
		RunE:   setRunE,
	}

	// Define flags
	SetCmd.Flags().StringP("username", "u", "", "Git username")
	SetCmd.Flags().StringP("remoteUrl", "r", "", "Remote repository URL")

	return SetCmd
}

func setRunE(cmd *cobra.Command, args []string) error{
	// Read flag values
	userName, _ := cmd.Flags().GetString("username")
	remoteUrl, _ := cmd.Flags().GetString("remoteUrl")

	configPath := filepath.Join(".rec", "config.json")

	// check if config.json exists
	if !utils.CheckFileExist(configPath){
		// create a empty config.json
		utils.CreateFile(configPath)
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
	if strings.TrimSpace(userName) != "" {
		config["userName"] = userName
	}
	if strings.TrimSpace(remoteUrl) != "" {
		config["remoteUrl"] = remoteUrl
	}

	// Write back to config.json
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error (set.go): ",err.Error())
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error (set.go): ",err.Error())
	}

	fmt.Println("✅ Repository configuration updated successfully!")

	return nil
}
