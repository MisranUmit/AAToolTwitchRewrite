package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Config struct {
	MinecraftPath string `json:"minecraftPath"`
}

func readConfigFile() Config {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	var decodedJson Config
	err = json.Unmarshal(configFile, &decodedJson)
	if err != nil {
		log.Fatal(err)
	}
	return decodedJson
}

func writeConfigFile(newConfig Config) {
	rawJson, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("config.json", rawJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func checkConfigValidity(config Config) bool {
	_, err := os.Stat(filepath.Join(config.MinecraftPath, "saves"))
	return err == nil
}

func getMostRecentWorld(savesPath string) (string, error) {
	dirList, err := os.ReadDir(savesPath)
	if err != nil {
		return "", err
	}

	type worldInfo struct {
		name    string
		modTime time.Time
		hasAdv  bool
	}

	var worlds []worldInfo

	for _, dir := range dirList {
		if !dir.IsDir() {
			continue
		}
		fileInfo, err := dir.Info()
		if err != nil {
			continue
		}

		// Check if this world has advancement data
		advPath := filepath.Join(savesPath, dir.Name(), "advancements")
		_, advErr := os.Stat(advPath)
		hasAdvancements := advErr == nil

		worlds = append(worlds, worldInfo{
			name:    dir.Name(),
			modTime: fileInfo.ModTime(),
			hasAdv:  hasAdvancements,
		})
	}

	if len(worlds) == 0 {
		return "", fmt.Errorf("no worlds found")
	}

	// Prefer worlds with advancement data, sorted by modification time
	var worldsWithAdv []worldInfo
	for _, w := range worlds {
		if w.hasAdv {
			worldsWithAdv = append(worldsWithAdv, w)
		}
	}

	var selectedWorld worldInfo
	if len(worldsWithAdv) > 0 {
		selectedWorld = worldsWithAdv[0]
		for _, w := range worldsWithAdv {
			if w.modTime.After(selectedWorld.modTime) {
				selectedWorld = w
			}
		}
	} else {
		// No worlds with advancement data, use the most recently modified
		selectedWorld = worlds[0]
		for _, w := range worlds {
			if w.modTime.After(selectedWorld.modTime) {
				selectedWorld = w
			}
		}
	}

	return selectedWorld.name, nil
}

func advancementToString(adv Advancement) string {
	return fmt.Sprintf("[%s] %s - %s", adv.Id, adv.Data.Display.Title, adv.Data.Display.Description)
}

func main() {
	configData := readConfigFile()

	if !checkConfigValidity(configData) {
		for {
			fmt.Print("Current path to .minecraft invalid, please enter new: ")

			reader := bufio.NewReader(os.Stdin)
			userInputPath, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			userInputPath = strings.TrimSpace(userInputPath)

			configData.MinecraftPath = userInputPath

			if checkConfigValidity(configData) {
				writeConfigFile(configData)
				break
			}
		}
	}

	savesPath := filepath.Join(configData.MinecraftPath, "saves")
	worldName, err := getMostRecentWorld(savesPath)
	if err != nil {
		log.Fatal(err)
	}

	worldPath := filepath.Join(savesPath, worldName)
	advancementsPath := filepath.Join(worldPath, "advancements")

	fmt.Printf("Monitoring world: %s\n", worldName)
	fmt.Printf("Watching: %s\n", advancementsPath)

	// Load initial advancements
	initialAdvancements, err := readAdvancementFiles(worldPath)
	if err != nil {
		log.Fatal(err)
	}

	seenAdvancements := make(map[string]bool)
	for _, adv := range initialAdvancements {
		seenAdvancements[adv.Id] = true
	}

	fmt.Printf("Loaded %d initial advancements\n\n", len(seenAdvancements))

	// Set up file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Watch the advancements directory and all subdirectories
	err = filepath.Walk(advancementsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Waiting for new advancements...")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				// Wait a moment for the file to be fully written
				time.Sleep(100 * time.Millisecond)

				currentAdvancements, err := readAdvancementFiles(worldPath)
				if err != nil {
					log.Printf("Error reading advancements: %v\n", err)
					continue
				}

				// Check for new advancements
				for _, adv := range currentAdvancements {
					if !seenAdvancements[adv.Id] {
						fmt.Printf("🎉 NEW ADVANCEMENT: %s\n", advancementToString(adv))
						seenAdvancements[adv.Id] = true
					}
				}
			}

		case event, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v\n", event)
		}
	}
}
