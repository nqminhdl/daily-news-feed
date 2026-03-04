package config

import (
	"os"

	"daily-news-feed/pkg/util"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Categories     map[string]Category `yaml:"categories"`
	FeedConfig     FeedConfig          `yaml:"feedConfig"`
	PositionConfig struct {
		Backend    string `yaml:"backend"`
		Filesystem struct {
			Path string `yaml:"path"`
		} `yaml:"filesystem"`
		Sqlite struct {
			Path     string `yaml:"path"`
			Database string `yaml:"database"`
		} `yaml:"sqlite"`
		MySQl struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"mysql"`
		PostgreSQL struct {
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"postgresql"`
		Redis struct {
			Enabled  bool   `yaml:"enabled"`
			Host     string `yaml:"host"`
			Port     string `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"redis"`
	} `yaml:"positionConfig"`
}

type FeedConfig struct {
	MaxAgeInDays int `yaml:"maxAgeInDays"` // Maximum age of feed items to fetch in days, defaults to 30 if not specified
}

type Category struct {
	Telegram struct {
		Enabled  bool   `yaml:"enabled"`
		ChatID   string `yaml:"chatId"`
		BotToken string `yaml:"botToken"`
	} `yaml:"telegram"`
	Prometheus struct {
		Enabled   bool   `yaml:"enabled"`
		Url       string `yaml:"url"`
		BasicAuth struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"baicAuth"`
	} `yaml:"prometheus"`
	Slack struct {
		Enabled    bool   `yaml:"enabled"`
		WebhookURL string `yaml:"webhookUrl"`
		BotToken   string `yaml:"botToken"`
		ChannelID  string `yaml:"channelId"`
	} `yaml:"slack"`
	Feed []Feed `yaml:"feeds"`
}

type Feed struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

func ReadConfig() Config {
	logger := util.Logger()
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		logger.Fatalf("error reading YAML file: %v", err)
	}

	expandedYAML := os.ExpandEnv(string(data))

	// Step 3: Unmarshal YAML into struct
	var config Config
	err = yaml.Unmarshal([]byte(expandedYAML), &config)
	if err != nil {
		logger.Fatalf("Error unmarshalling YAML: %s\n", err)
	}

	return config
}
