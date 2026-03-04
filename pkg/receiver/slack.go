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

const slackChatPostMessageURL = "https://slack.com/api/chat.postMessage"

type slackChatPostMessageRequest struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type slackChatPostMessageResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

func sendSlackMessage(webhookURL string, botToken string, channelID string, title string, url string) error {
	logger := util.Logger()

	message := fmt.Sprintf("%s\n%s", title, url)
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	if botToken != "" || channelID != "" {
		if botToken == "" || channelID == "" {
			logger.Error("Slack bot token and channel ID must both be set to use Slack API")
		} else {
			if err := sendSlackMessageByAPI(client, botToken, channelID, message); err == nil {
				logger.Debugf("Message sent to Slack channel successfully: %s", title)
				return nil
			} else {
				logger.Errorf("Failed to send message using Slack API: %v", err)
				if webhookURL == "" {
					return err
				}
				logger.Warn("Falling back to Slack incoming webhook")
			}
		}
	}

	if webhookURL == "" {
		logger.Error("Slack webhook URL is not set")
		return nil
	}

	if err := sendSlackMessageByWebhook(client, webhookURL, message); err != nil {
		logger.Errorf("Failed to send message using Slack webhook: %v", err)
		return err
	}

	logger.Debugf("Message sent to Slack successfully: %s", title)
	return nil
}

func sendSlackMessageByWebhook(client *http.Client, webhookURL string, message string) error {
	payload := map[string]string{"text": message}
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal Slack webhook payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("create Slack webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send Slack webhook request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("unexpected Slack webhook status code %d: %s", res.StatusCode, body)
	}

	return nil
}

func sendSlackMessageByAPI(client *http.Client, botToken string, channelID string, message string) error {
	payload := slackChatPostMessageRequest{
		Channel: channelID,
		Text:    message,
	}
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal Slack API payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, slackChatPostMessageURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("create Slack API request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", botToken))

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send Slack API request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read Slack API response: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected Slack API status code %d: %s", res.StatusCode, body)
	}

	var response slackChatPostMessageResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("unmarshal Slack API response: %w", err)
	}

	if !response.OK {
		if response.Error == "" {
			response.Error = "unknown_error"
		}
		return fmt.Errorf("Slack API error: %s", response.Error)
	}

	return nil
}
