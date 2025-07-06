package backend

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"daily-news-feed/pkg/util"
)

type FSData struct {
	Name    string `yaml:"name"`
	Link    string `yaml:"link"`
	PubDate int64  `yaml:"pubDate"`
}

type FSPositionData struct {
	Positions []FSData `yaml:"positions"`
}

func verifyYaml(filename string) (*FSPositionData, error) {
	logger := util.Logger()
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &FSPositionData{Positions: []FSData{}}, nil
		}
		logger.Errorf("Failed to read YAML file: %v", err)
		return nil, err
	}

	var fsPositionData FSPositionData
	if err := yaml.Unmarshal(data, &fsPositionData); err != nil {
		logger.Errorf("Failed to unmarshal YAML: %v", err)
		return nil, err
	}

	if fsPositionData.Positions == nil {
		fsPositionData.Positions = []FSData{}
	}

	return &fsPositionData, nil
}

func FsDataWriting(filename string, name string, link string, pubDate string) (bool, error) {
	logger := util.Logger()
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("Failed to open or create YAML file: %v", err)
		return false, err
	}
	defer file.Close()

	positionConfig, err := verifyYaml(filename)
	if err != nil {
		return false, err
	}

	// Check if link already exists
	for _, position := range positionConfig.Positions {
		if position.Link == link {
			logger.Debugf("Position '%s' already exists", position.Name)
			return true, nil
		}
	}

	// Add new position
	pubDateInt, err := strconv.ParseInt(pubDate, 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse pubDate: %v", err)
		return false, err
	}

	newPosition := FSData{
		Name:    name,
		Link:    link,
		PubDate: pubDateInt,
	}

	logger.Debugf("Adding new position '%s'", newPosition.Name)
	positionConfig.Positions = append(positionConfig.Positions, newPosition)

	// Write updated data
	updateYaml, err := yaml.Marshal(positionConfig)
	if err != nil {
		logger.Errorf("Failed to marshal YAML: %v", err)
		return false, err
	}

	if err := os.WriteFile(filename, updateYaml, 0644); err != nil {
		logger.Errorf("Failed to write YAML file: %v", err)
		return false, err
	}

	return false, nil // New item added
}
