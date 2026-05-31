package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type AdvancementData struct {
	Display struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Icon        struct {
			Item string `json:"item"`
		} `json:"icon"`
		Frame      string `json:"frame"`
		ShowToast  bool   `json:"show_toast"`
		AnnounceToChat bool   `json:"announce_to_chat"`
	} `json:"display"`
	Criteria map[string]interface{} `json:"criteria"`
	Parent   string                 `json:"parent"`
}

type Advancement struct {
	Id   string
	Data AdvancementData
}

// readAdvancementFiles reads all advancements from a world's advancements directory
func readAdvancementFiles(worldPath string) ([]Advancement, error) {
	advancementsPath := filepath.Join(worldPath, "advancements")
	
	// Check if advancements directory exists
	if _, err := os.Stat(advancementsPath); err != nil {
		return nil, fmt.Errorf("advancements directory not found: %v", err)
	}

	var advancements []Advancement

	// Walk through all JSON files in the advancements directory
	err := filepath.Walk(advancementsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".json") {
			// Read and parse the JSON file
			fileContent, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading advancement file %s: %v\n", path, err)
				return nil
			}

			var advData AdvancementData
			err = json.Unmarshal(fileContent, &advData)
			if err != nil {
				log.Printf("Error parsing advancement file %s: %v\n", path, err)
				return nil
			}

			// Extract the advancement ID from the file path
			// Path format: <world>/advancements/minecraft/story/mine_wood.json -> minecraft:story/mine_wood
			relPath, _ := filepath.Rel(advancementsPath, path)
			advID := strings.TrimSuffix(relPath, ".json")
			advID = strings.ReplaceAll(advID, "\\", "/")
			advID = strings.Replace(advID, "/", ":", 1)

			advancements = append(advancements, Advancement{
				Id:   advID,
				Data: advData,
			})
		}

		return nil
	})

	return advancements, err
}
