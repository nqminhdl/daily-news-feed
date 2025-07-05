package receiver

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"daily-news-feed/pkg/util"
)

func sendTelegramMessage(botToken string, chatId string, url string) error {
	logger := util.Logger()
	if botToken == "" {
		logger.Error("Telegram bot token is not set")
		return nil
	}

	if chatId == "" {
		logger.Error("Telegram chat ID is not defined")
		return nil
	}

	baseURL := os.Getenv("TELEGRAM_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.telegram.org/bot"
	}

	sendMessageURL := fmt.Sprintf("%s%s/sendMessage?chat_id=%s&text=%s", baseURL, botToken, chatId, url)

	req, err := http.NewRequest("GET", sendMessageURL, nil)
	if err != nil {
		logger.Errorf("Failed to create request: %v", err)
		return err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := client.Do(req)
	if err != nil {
		logger.Errorf("Failed to send request: %v", err)
		return err
	}
	defer response.Body.Close()

	// Sleep 2 seconds to avoid 429 error from Telegram
	time.Sleep(2 * time.Second)

	if response.StatusCode != http.StatusOK {
		logger.Errorf("Unexpected status code: %d", response.StatusCode)
		return err
	}

	logger.Debugf("Message sent to Telegram successfully: %s", url)
	return nil
}
