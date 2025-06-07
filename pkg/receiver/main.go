package receiver

import (
	"daily-news-feed/pkg/config"

	util "daily-news-feed/pkg/util"
)

func SendNotification(config *config.Category, category string, title string, link string, pubDate string) {
	logger := util.Logger()

	if config.Prometheus.Enabled {
		logger.Infof("Prometheus is enabled, produce metrics to %s.\n", config.Prometheus.Url)
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
		logger.Infof("Telegram is enabled, send %s to channel ID%s.\n", link, config.Telegram.ChatID)
		sendTelegramMessage(
			config.Telegram.BotToken,
			config.Telegram.ChatID,
			link,
		)
	}

	if config.Slack.Enabled {
		logger.Infof("Slack is enabled, sending %s.\n", link)
		sendSlackMessage(
			config.Slack.WebhookUrlUrl,
			title,
			link,
		)
	}
}
