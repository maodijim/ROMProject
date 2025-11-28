package main

import (
	"fmt"
	"os"

	"ROMProject/tools/webApp/usersSpace"

	"gopkg.in/yaml.v3"
)

// Add helper functions here if needed, e.g., for reading/writing config files.
func saveConfigs() error {
	updatedData, err := yaml.Marshal(&usersSpace.Configs)
	if err != nil {
		return fmt.Errorf("failed to marshal configs: %w", err)
	}

	if err := os.WriteFile(configPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
