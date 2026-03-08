package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	// Future proofing in case I ever add more options
	MinecraftPath string `json:"minecraftPath"`
}

// TODO: Make this function return err like the others
func readConfigFile() Config {
	// Retrieve file to raw text in RAM
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Decode json?
	var decodedJson Config
	err = json.Unmarshal(configFile, &decodedJson)
	if err != nil {
		log.Fatal(err)
	}
	return decodedJson
}

func writeConfigFile(newConfig Config) {
	// Turns the struct back into raw json data
	rawJson, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Dunno if 0644 is good practice, but it's what the interwebz says
	err = os.WriteFile("config.json", rawJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

/*
func readAdvancementFiles() {
	advancementFile, err := os.ReadFile(".\\advancementAssets\\minecraft.xml")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(advancementFile))
}
*/

// Checks if saves folder is present withing path, if user manages to give path that isn't a
// Minecraft folder but simultaneously has a saves folder then good luck to them.
func checkConfigValidity(config Config) bool {
	_, err := os.Stat(config.MinecraftPath + "\\saves")
	if err != nil {
		return false
	}
	return true
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
	dirList, err := os.ReadDir(configData.MinecraftPath + "\\saves")
	if err != nil {
		log.Fatal(err)
	}
	var newest time.Time
	fmt.Println(newest)
	for i := range len(dirList) {
		fileInfo, err := dirList[i].Info()
		if err != nil {
			log.Fatal(err)
		}
		if newest.Compare(fileInfo.ModTime()) == -1 {
			newest = fileInfo.ModTime()
		} else if newest.Compare(fileInfo.ModTime()) == 0 {
			log.Fatal("Two worlds are created at the same time, unable to gauge which one to choose.")
		}
	}
	fmt.Println(newest)
	fuck := readAdvancementFiles()
	fmt.Println(fuck.Group.Id)
}
