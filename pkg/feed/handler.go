package handler

import (
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"

	"daily-news-feed/pkg/backend"
	"daily-news-feed/pkg/config"
	"daily-news-feed/pkg/receiver"
	"daily-news-feed/pkg/util"
)

func FeedHandler() {
	logger := util.Logger()
	fp := gofeed.NewParser()
	fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	cfg := config.ReadConfig()
	categories := cfg.Categories
	backendConfig := cfg.PositionConfig

	// Get max age from config, default to 30 days if not specified
	maxAge := 30
	if cfg.FeedConfig.MaxAgeInDays > 0 {
		maxAge = cfg.FeedConfig.MaxAgeInDays
	}
	cutoffTime := time.Now().AddDate(0, 0, -maxAge)

	for name, config := range categories {
		for _, feed := range config.Feed {
			logger.Infof("Parsing feed - Category: %s, Feed: %s", name, feed.Name)

			parsedURL, err := fp.ParseURL(feed.URL)
			if err != nil {
				logger.Errorf("Failed to parse URL %s: %v", feed.URL, err)
				continue // Skip this feed but continue with others
			}

			for _, item := range parsedURL.Items {
				// Skip items older than maxAge
				if item.PublishedParsed != nil && item.PublishedParsed.Before(cutoffTime) {
					logger.Debugf("Skipping old item: %s (published: %s)", item.Title, item.PublishedParsed)
					continue
				}

				logger.Debugf("Processing item: %s", item.Title)
				pubDate := getPubDate(item)

				exists, err := storeBackend(backendConfig.Backend, item.Title, item.Link, pubDate, backendConfig.Filesystem.Path, backendConfig.Sqlite.Path)
				if err != nil {
					logger.Errorf("Failed to store item %s: %v", item.Title, err)
					continue
				}

				// Only send notification for new items
				if !exists {
					receiver.SendNotification(&config, name, item.Title, item.Link, pubDate)
				}
			}
		}
	}
}

// getPubDate returns the publication date as a Unix timestamp string
func getPubDate(item *gofeed.Item) string {
	if item.PublishedParsed != nil {
		return strconv.FormatInt(item.PublishedParsed.Unix(), 10)
	}
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// storeBackend stores an item in the configured backend
func storeBackend(backendType, title, link, pubDate, fsPath, sqlitePath string) (bool, error) {
	logger := util.Logger()
	switch backendType {
	case "filesystem":
		return backend.FsDataWriting(fsPath, title, link, pubDate)
	case "sqlite":
		return backend.SQLiteWriting(sqlitePath, title, link, pubDate)
	default:
		logger.Errorf("unsupported backend type: %s", backendType)
		return false, nil
	}
}
