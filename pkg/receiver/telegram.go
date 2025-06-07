package receiver

import (
	"fmt"
	"net/http"
	"os"
	"time"

	util "daily-news-feed/pkg/util"
)

func sendTelegramMessage(botToken string, chatId string, url string) error {
	logger := util.Logger()
	if botToken == "" {
		return fmt.Errorf("telegram bot token is not set")
	}

	if chatId == "" {
		return fmt.Errorf("telegram chat id is not defined")
	}

	baseURL := os.Getenv("TELEGRAM_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.telegram.org/bot"
	}

	sendMessageURL := fmt.Sprintf("%s%s/sendMessage?chat_id=%s&text=%s", baseURL, botToken, chatId, url)

	req, err := http.NewRequest("GET", sendMessageURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Sleep 2 seconds to avoid 429 error from Telegram
	time.Sleep(2 * time.Second)

	if response.StatusCode != http.StatusOK {
		logger.Infof("unexpected status code: %d", response.StatusCode)
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}
