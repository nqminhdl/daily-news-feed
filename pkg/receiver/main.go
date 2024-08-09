package receiver

import (
	"daily-news-feed/pkg/config"
	"log"
)

func SendNotification(config *config.Category, category string, title string, link string, pubDate string) {
	if config.Prometheus.Enabled {
		log.Printf("Prometheus is enabled, produce metrics to %s.\n", config.Prometheus.Url)
		produceMetricsToPrometheus(
			config.Prometheus.BasicAuth.Username,
			config.Prometheus.BasicAuth.Password,
			config.Prometheus.Url,
			category,
			title,
			link,
			pubDate,
		)
	}

	if config.Telegram.Enabled {
		log.Printf("Telegram is enabled, send %s to channel ID%s.\n", link, config.Telegram.ChatID)
		sendTelegramMessage(
			config.Telegram.BotToken,
			config.Telegram.ChatID,
			link,
		)
	}

	if config.Slack.Enabled {
		log.Printf("Slack is enabled, sending %s.\n", link)
		sendSlackMessage(
			config.Slack.WebhookUrlUrl,
			title,
			link,
		)
	}
}
