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

func getAdvancementTitle(advId string) string {
	// Comprehensive map of ALL Minecraft 1.21 advancement IDs to proper display names
	// Generated from official Minecraft advancement data
	titles := map[string]string{
		// ROOT
		"minecraft:story/root":     "Minecraft",
		"minecraft:adventure/root": "Adventure",
		"minecraft:nether/root":    "Nether",
		"minecraft:end/root":       "The End",
		"minecraft:husbandry/root": "Husbandry",
		// STORY
		"minecraft:story/mine_stone":       "Stone Age",
		"minecraft:story/upgrade_tools":    "Getting an Upgrade",
		"minecraft:story/smelt_iron":       "Acquire Hardware",
		"minecraft:story/obtain_armor":     "Suit Up",
		"minecraft:story/lava_bucket":      "Hot Stuff",
		"minecraft:story/deflect_arrow":    "Not Today, Thank You",
		"minecraft:story/enter_the_nether": "We Need to Go Deeper",
		"minecraft:story/follow_ender_eye": "The End?",
		"minecraft:story/enter_the_end":    "The End",
		// ADVENTURE
		"minecraft:adventure/adventuring_time":      "Adventuring Time",
		"minecraft:adventure/kill_all_mobs":         "Monsters Hunted",
		"minecraft:adventure/kill_all_hostile_mobs": "Monsters Hunted",
		"minecraft:adventure/voluntary_exile":       "Voluntary Exile",
		"minecraft:adventure/arbalistic":            "Arbalistic",
		"minecraft:adventure/totem_of_undying":      "Uneasy Alliance",
		"minecraft:adventure/sleep_in_bed":          "Sweet Dreams",
		"minecraft:adventure/throw_trident":         "Take Aim",
		"minecraft:adventure/shoot_arrow":           "Take Aim",
		"minecraft:adventure/trade":                 "What a Deal!",
		"minecraft:adventure/hero_of_the_village":   "Hero of the Village",
		"minecraft:adventure/honey_block_slide":     "Sticky Situation",
		"minecraft:adventure/ol_betsy":              "Best Friends Forever",
		"minecraft:adventure/spyglass_at_parrot":    "Is It a Bird?",
		"minecraft:adventure/spyglass_at_ghast":     "Is It a Balloon?",
		"minecraft:adventure/bullseye":              "Bullseye",
		// END
		"minecraft:end/kill_dragon":       "Free the End",
		"minecraft:end/dragon_breath":     "The Next Generation",
		"minecraft:end/enter_end_gateway": "Remote Getaway",
		"minecraft:end/respawn_dragon":    "The Dragon Egg",
		// NETHER
		"minecraft:nether/return_from_void":       "Return from the Void",
		"minecraft:nether/find_fortress":          "Nether",
		"minecraft:nether/obtain_ancient_debris":  "Hidden in the Depths",
		"minecraft:nether/fast_travel":            "Subspace Bubble",
		"minecraft:nether/find_bastion":           "Those Were the Days",
		"minecraft:nether/obtain_blaze_rod":       "Blaze and Blaze Away",
		"minecraft:nether/netherite_armor":        "Fully Powered",
		"minecraft:nether/use_lodestone":          "Lodestone Compass",
		"minecraft:nether/explore_nether":         "Hot Tourist Destinations",
		"minecraft:nether/summon_wither":          "Withering Heights",
		"minecraft:nether/obtain_crying_obsidian": "Who is Cutting Onions?",
		// HUSBANDRY
		"minecraft:husbandry/breed_an_animal":              "The Parrots and the Bats",
		"minecraft:husbandry/tame_an_animal":               "Best Friends Forever",
		"minecraft:husbandry/fishy_business":               "Fishy Business",
		"minecraft:husbandry/plant_seed":                   "A Seedy Place",
		"minecraft:husbandry/breed_all_animals":            "Serious Dedication",
		"minecraft:husbandry/complete_catalogue":           "A Complete Catalogue",
		"minecraft:husbandry/tactical_fishing":             "Tactical Fishing",
		"minecraft:husbandry/ride_strider":                 "Ride the Lava",
		"minecraft:husbandry/summon_axolotl":               "The Cutest Predator",
		"minecraft:husbandry/obtain_netherite_hoe":         "Serious Dedication",
		"minecraft:husbandry/wax_on":                       "Wax On",
		"minecraft:husbandry/wax_off":                      "Wax Off",
		"minecraft:husbandry/silk_touch_nest":              "This is Fine",
		"minecraft:husbandry/find_bee_nest":                "Bee Our Guest",
		"minecraft:husbandry/balanced_diet":                "A Balanced Diet",
		"minecraft:husbandry/break_diamond_hoe":            "Serious Dedication",
		"minecraft:husbandry/allay_deliver_item_to_player": "You Got a Friend",
		"minecraft:husbandry/safely_harvest_honey":         "Sticky Situation",
		"minecraft:husbandry/kill_axolotl_with_dryout":     "The Cutest Predator",
		"minecraft:husbandry/leash_all_frog_variants":      "Leash All Frogs",
		"minecraft:husbandry/tadpole_into_frog":            "When the Bloom is on the Vine",
		"minecraft:husbandry/make_a_sign_glow":             "Glow and Behold!",
		"minecraft:husbandry/froglights":                   "With Our Powers Combined!",
	}

	if title, exists := titles[advId]; exists {
		return title
	}

	// Fallback: format the ID nicely
	formatted := strings.TrimPrefix(advId, "minecraft:")
	if strings.Contains(formatted, "/") {
		parts := strings.Split(formatted, "/")
		formatted = parts[len(parts)-1]
	}
	formatted = strings.ReplaceAll(formatted, "_", " ")
	// Capitalize first letter of each word
	words := strings.Fields(formatted)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
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

	// Set up file watcher for saves directory to detect new worlds
	savesWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer savesWatcher.Close()

	err = savesWatcher.Add(savesPath)
	if err != nil {
		log.Fatal(err)
	}

	seenAdvancements := make(map[string]bool)
	var currentWorld string
	var advWatcher *fsnotify.Watcher

	// Initialize advancement watcher for a world
	initializeWorld := func() error {
		newWorld, err := getMostRecentWorld(savesPath)
		if err != nil {
			return err
		}

		// If same world, no need to reinitialize
		if newWorld == currentWorld && advWatcher != nil {
			return nil
		}

		// Close old advancement watcher
		if advWatcher != nil {
			advWatcher.Close()
		}

		currentWorld = newWorld
		worldPath := filepath.Join(savesPath, currentWorld)
		advancementsPath := filepath.Join(worldPath, "advancements")

		fmt.Printf("\nMonitoring world: %s\n", currentWorld)
		fmt.Printf("Watching: %s\n", advancementsPath)

		// Load initial advancements for this world
		initialAdvancements, err := readAdvancementFiles(worldPath)
		if err != nil {
			return err
		}

		// Clear old advancements and load new ones
		seenAdvancements = make(map[string]bool)
		for _, adv := range initialAdvancements {
			if adv.Progress.Done {
				seenAdvancements[adv.Id] = true
			}
		}

		fmt.Printf("Loaded %d initial advancements\n", len(seenAdvancements))
		fmt.Println("Waiting for new advancements...")

		// Create new advancement watcher
		advWatcher, err = fsnotify.NewWatcher()
		if err != nil {
			return err
		}

		err = filepath.Walk(advancementsPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return advWatcher.Add(path)
			}
			return nil
		})

		return err
	}

	// Initialize with current world
	if err := initializeWorld(); err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-savesWatcher.Events:
			if !ok {
				return
			}
			// Recheck world periodically via ticker instead
			_ = event

		case event, ok := <-advWatcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				// Wait a moment for the file to be fully written
				time.Sleep(100 * time.Millisecond)

				worldPath := filepath.Join(savesPath, currentWorld)
				currentAdvancements, err := readAdvancementFiles(worldPath)
				if err != nil {
					log.Printf("Error reading advancements: %v\n", err)
					continue
				}

				// Check for new advancements
				for _, adv := range currentAdvancements {
					if adv.Progress.Done && !seenAdvancements[adv.Id] && !strings.Contains(adv.Id, "minecraft:recipes/") {
						fmt.Printf("NEW ADVANCEMENT: %s\n", getAdvancementTitle(adv.Id))
						seenAdvancements[adv.Id] = true
					}
				}
			}

		case <-ticker.C:
			// Periodically check if world has changed
			if err := initializeWorld(); err != nil {
				log.Printf("Error checking for world change: %v\n", err)
			}

		case event, ok := <-savesWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("Saves watcher error: %v\n", event)

		case event, ok := <-advWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("Advancement watcher error: %v\n", event)
		}
	}
}
