package receiver

import (
	"daily-news-feed/pkg/config"

	"daily-news-feed/pkg/util"
)

func SendNotification(config *config.Category, category string, title string, link string, pubDate string) {
	logger := util.Logger()

	if config.Prometheus.Enabled {
		logger.Debugf("Prometheus is enabled, sending metrics to %s.\n", config.Prometheus.Url)
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
		logger.Debugf("Telegram is enabled, send %s to channel ID%s.\n", link, config.Telegram.ChatID)
		sendTelegramMessage(
			config.Telegram.BotToken,
			config.Telegram.ChatID,
			link,
		)
	}

	if config.Slack.Enabled {
		logger.Debugf("Slack is enabled, sending %s.\n", link)
		sendSlackMessage(
			config.Slack.WebhookURL,
			config.Slack.BotToken,
			config.Slack.ChannelID,
			title,
			link,
		)
	}
}
