package handler

import (
	"log"
	"strconv"

	backend "daily-news-feed/pkg/backend"
	config "daily-news-feed/pkg/config"
	receiver "daily-news-feed/pkg/receiver"

	"github.com/mmcdole/gofeed"
)

func FeedHandler() {
	fp := gofeed.NewParser()
	fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	categories := config.ReadConfig().Categories
	backendConfig := config.ReadConfig().PositionConfig

	for name, config := range categories {
		for _, feed := range config.Feed {
			log.Printf("Parsing Links - Category %s - %s", name, feed.Name)

			parsedURL, err := fp.ParseURL(feed.URL)
			if err != nil {
				log.Fatalf("error parsing URL: %v", err)
			}

			for _, item := range parsedURL.Items {
				log.Printf("Checking backend for item: %s", item.Title)
				pubDate := strconv.FormatInt(item.PublishedParsed.Unix(), 10)
				switch backendConfig.Backend {
				case `filesystem`:
					linkFound := backend.FsDataWriting(backendConfig.Filesystem.Path, item.Title, item.Link, pubDate)
					if !linkFound {
						receiver.SendNotification(&config, name, item.Title, item.Link, pubDate)
					}
				case `sqlite`:
					linkFound, err := backend.SQLiteWriting(backendConfig.Sqlite.Path, item.Title, item.Link, pubDate)
					if !linkFound && err == nil {
						receiver.SendNotification(&config, name, item.Title, item.Link, pubDate)
					}
				default:
					log.Printf("Unhandled backend: %s", backendConfig.Backend)
				}
			}
		}
	}
}
