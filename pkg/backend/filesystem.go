package backend

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type FSData struct {
	Name    string `yaml:"name"`
	Link    string `yaml:"link"`
	PubDate int64  `yaml:"pubDate"`
}

type FSPositionData struct {
	Positions []FSData `yaml:"positions"`
}

func verfifyYaml(filename string) FSPositionData {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading YAML file: %v", err)
	}

	var fsPositionData FSPositionData
	err = yaml.Unmarshal(data, &fsPositionData)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
	}
	return fsPositionData
}

func FsDataWriting(filename string, name string, link string, pubDate string) bool {
	fileName := filename

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("error opening or creating YAML file: %v", err)
	}
	defer file.Close()

	positionConfig := verfifyYaml(fileName)

	linkFound := false
	for _, position := range positionConfig.Positions {
		if position.Link == link {
			log.Printf("Position '%s' is already in the list", position.Name)
			linkFound = true
			break
		}
	}

	if linkFound {
		return linkFound
	}

	pubDateInt, _ := strconv.ParseInt(pubDate, 10, 32)
	newPosition := FSData{
		Name:    name,
		Link:    link,
		PubDate: pubDateInt,
	}

	log.Printf("Position '%s' is not in the list. Writing date to position file.", newPosition.Name)
	positionConfig.Positions = append(positionConfig.Positions, newPosition)

	updateYaml, err := yaml.Marshal(&positionConfig)
	if err != nil {
		log.Fatalf("error marshalling YAML: %v", err)
	}

	err = os.WriteFile(fileName, updateYaml, 0644)
	if err != nil {
		log.Fatalf("error writing YAML file: %v", err)
	}

	return linkFound
}
