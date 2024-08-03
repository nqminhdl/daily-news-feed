package backend

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteData struct {
	ID      int64
	Name    string
	Link    string
	PubDate string
}

func SQLiteWriting(filename string, name string, link string, pubDate string) (bool, error) {
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&SQLiteData{}); err != nil {
		panic("failed to migrate database")
	}

	var data SQLiteData
	linkFound := false

	newPosition := SQLiteData{
		Name:    name,
		Link:    link,
		PubDate: pubDate,
	}
	result := db.Where("link = ?", link).First(&data)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&newPosition).Error; err != nil {
				log.Printf("Failed to insert new position '%s': %v", link, err)
				return linkFound, err
			}
			log.Printf("New position '%s' added to the database", link)
		} else {
			log.Printf("Error querying the database for position '%s': %v", link, result.Error)
			return linkFound, result.Error
		}
	} else {
		log.Printf("Position '%s' is already in the list", link)
		linkFound = true
	}

	return linkFound, nil
}
