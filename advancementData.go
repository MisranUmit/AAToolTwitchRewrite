package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AdvancementProgress struct {
	Criteria map[string]interface{} `json:"criteria"`
	Done     bool                   `json:"done"`
}

type Advancement struct {
	Id       string
	Progress AdvancementProgress
}

// readAdvancementFiles reads all advancements for all players in a world
// Minecraft stores advancements per-player as UUID.json files
func readAdvancementFiles(worldPath string) ([]Advancement, error) {
	advancementsPath := filepath.Join(worldPath, "advancements")
	
	// Check if advancements directory exists
	if _, err := os.Stat(advancementsPath); err != nil {
		return nil, fmt.Errorf("advancements directory not found: %v", err)
	}

	var allAdvancements []Advancement

	// Read all UUID-based advancement files
	files, err := os.ReadDir(advancementsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading advancements directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(advancementsPath, file.Name())
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading advancement file %s: %v\n", filePath, err)
			continue
		}

		// Parse the player's advancement file
		// Format: { "advancement/id": { "criteria": {...}, "done": true/false }, ... }
		var playerAdvancements map[string]interface{}
		err = json.Unmarshal(fileContent, &playerAdvancements)
		if err != nil {
			log.Printf("Error parsing advancement file %s: %v\n", filePath, err)
			continue
		}

		// Convert to our advancement list
		for advId, advDataRaw := range playerAdvancements {
			// Skip metadata fields
			if advId == "DataVersion" {
				continue
			}

			// Parse the advancement object
			advBytes, _ := json.Marshal(advDataRaw)
			var progress AdvancementProgress
			err := json.Unmarshal(advBytes, &progress)
			if err != nil {
				// Skip entries that don't have the advancement structure
				continue
			}

			allAdvancements = append(allAdvancements, Advancement{
				Id:       advId,
				Progress: progress,
			})
		}
	}

	return allAdvancements, nil
}
