package receiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"daily-news-feed/pkg/util"
)

func sendSlackMessage(webhookURL string, title string, url string) error {
	logger := util.Logger()
	if webhookURL == "" {
		logger.Error("Slack webhook URL is not set")
		return nil
	}

	message := fmt.Sprintf("%s\n%s", title, url)
	payload := map[string]string{"text": message}

	jsonStr, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("Failed to marshal payload: %v", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Errorf("Failed to create request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		logger.Errorf("Failed to send request: %v", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		logger.Errorf("Unexpected status code %d: %s", res.StatusCode, body)
		return err
	}

	logger.Debugf("Message sent to Slack successfully: %s", title)
	return nil
}
