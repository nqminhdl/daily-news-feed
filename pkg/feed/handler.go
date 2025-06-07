package handler

import (
	"strconv"
	"time"

	util "daily-news-feed/pkg/util"

	backend "daily-news-feed/pkg/backend"
	config "daily-news-feed/pkg/config"
	receiver "daily-news-feed/pkg/receiver"

	"github.com/mmcdole/gofeed"
)

func FeedHandler() {
	logger := util.Logger()
	fp := gofeed.NewParser()
	fp.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"

	categories := config.ReadConfig().Categories
	backendConfig := config.ReadConfig().PositionConfig

	for name, config := range categories {
		for _, feed := range config.Feed {
			logger.Infof("Parsing Links - Category %s - %s", name, feed.Name)

			parsedURL, err := fp.ParseURL(feed.URL)
			if err != nil {
				logger.Fatalf("error parsing URL: %v", err)
			}

			for _, item := range parsedURL.Items {
				logger.Infof("Checking backend for item: %s", item.Title)
				var pubDate string
				if item.PublishedParsed != nil {
					pubDate = strconv.FormatInt(item.PublishedParsed.Unix(), 10)
				} else {
					pubDate = strconv.FormatInt(time.Now().Unix(), 10)
				}
				switch backendConfig.Backend {
				case `filesystem`:
					receiver.SendNotification(&config, name, item.Title, item.Link, pubDate)
					backend.FsDataWriting(backendConfig.Filesystem.Path, item.Title, item.Link, pubDate)
				case `sqlite`:
					backend.SQLiteWriting(backendConfig.Sqlite.Path, item.Title, item.Link, pubDate)
					receiver.SendNotification(&config, name, item.Title, item.Link, pubDate)
				default:
					logger.Errorf("Unhandled backend: %s", backendConfig.Backend)
				}
			}
		}
	}
}
