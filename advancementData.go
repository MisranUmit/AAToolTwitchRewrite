package main

import (
	"encoding/xml"
	"log"
	"os"
)

/*
type AdvancementList struct {
	Advancement Advancement `xml:"group>group>advancement"`
}


type Advancement struct {
	Critera
}
*/

type Root struct {
	Group Advancement `xml:"group>advancement"`
}

type Advancement struct {
	Id string `xml:"id,attr"`
}

func readAdvancementFiles() Root {
	advancementFile, err := os.ReadFile(".\\advancementAssets\\minecraft.xml")
	if err != nil {
		log.Fatal(err)
	}

	var decodedXml Root
	err = xml.Unmarshal(advancementFile, &decodedXml)
	if err != nil {
		log.Fatal(err)
	}
	return decodedXml
}
