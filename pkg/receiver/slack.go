package receiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	util "daily-news-feed/pkg/util"
)

func sendSlackMessage(slackWebhookUrl string, title string, url string) error {
	logger := util.Logger()
	message := fmt.Sprintf("%s\n%s", title, url)

	payload := map[string]string{"text": message}
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		logger.Fatalf("client: could not marshal payload: %s\n", err)
	}

	req, err := http.NewRequest(http.MethodPost, slackWebhookUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Fatalf("client: could not create request: %s\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("fail to send message to Slack. Reponse code %s. Response message: %s", strconv.Itoa(res.StatusCode), res.Body)
	}

	logger.Info("Message successfully sent to Slack")
	return nil
}
