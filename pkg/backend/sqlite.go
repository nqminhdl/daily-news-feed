package backend

import (
	"daily-news-feed/pkg/util"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteData struct {
	ID      int64  `gorm:"primaryKey"`
	Name    string `gorm:"not null"`
	Link    string `gorm:"uniqueIndex;not null"`
	PubDate string `gorm:"not null"`
}

func SQLiteWriting(filename string, name string, link string, pubDate string) (bool, error) {
	logger := util.Logger()
	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return false, err
	}

	if err := db.AutoMigrate(&SQLiteData{}); err != nil {
		logger.Errorf("Failed to migrate database schema: %v", err)
		return false, err
	}

	// Check if link already exists
	var data SQLiteData
	result := db.Where("link = ?", link).First(&data)
	if result.Error == nil {
		logger.Debugf("Position '%s' already exists", name)
		return true, nil // Item already exists
	}
	if result.Error != gorm.ErrRecordNotFound {
		logger.Errorf("Failed to query database: %v", result.Error)
		return false, result.Error
	}

	// Create new record
	newPosition := SQLiteData{
		Name:    name,
		Link:    link,
		PubDate: pubDate,
	}
	if err := db.Create(&newPosition).Error; err != nil {
		logger.Errorf("Failed to insert new position '%s': %v", name, err)
		return false, err
	}

	logger.Debugf("Added new position '%s'", name)
	return false, nil // New item added
}
