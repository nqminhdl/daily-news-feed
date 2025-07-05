package receiver

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"time"

	"daily-news-feed/pkg/util"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

type RemoteWriteHandler struct {
	url     string
	headers map[string]string
	client  *http.Client
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (r *RemoteWriteHandler) WriteMetrics(data []byte) error {
	logger := util.Logger()
	req, err := http.NewRequest("POST", r.url, bytes.NewBuffer(data))
	if err != nil {
		logger.Errorf("Failed to create request: %v", err)
		return err
	}
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	time.Sleep(1 * time.Second)
	resp, err := r.client.Do(req)
	if err != nil {
		logger.Errorf("Failed to send request: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Unexpected status code: %s", resp.Status)
		return err
	}

	return nil
}

func createTimeSeries(metricName string, value float64, labels map[string]string) []prompb.TimeSeries {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	labelPairs := make([]prompb.Label, 0, len(labels)+1)
	labelPairs = append(labelPairs, prompb.Label{
		Name:  "__name__",
		Value: metricName,
	})
	for k, v := range labels {
		labelPairs = append(labelPairs, prompb.Label{
			Name:  k,
			Value: v,
		})
	}

	return []prompb.TimeSeries{
		{
			Labels: labelPairs,
			Samples: []prompb.Sample{
				{
					Value:     value,
					Timestamp: timestamp,
				},
			},
		},
	}
}

func produceMetricsToPrometheus(username string, password string, prometheusURL string, category string, title string, url string, pubDate string) error {
	logger := util.Logger()
	if prometheusURL == "" {
		logger.Error("Prometheus URL is not set")
		return nil
	}

	headers := map[string]string{
		"Authorization":                     "Basic " + basicAuth(username, password),
		"Content-Type":                      "application/x-protobuf",
		"X-Prometheus-Remote-Write-Version": "0.1.0",
	}

	labels := map[string]string{
		"category":     category,
		"service":      "daily-news-feed",
		"title":        title,
		"url":          url,
		"publish_date": pubDate,
	}

	handler := &RemoteWriteHandler{
		url:     prometheusURL,
		headers: headers,
		client:  &http.Client{Timeout: 10 * time.Second},
	}

	// Metric name feed_news_by_category always has a value of 1
	// This metric is simply hardcoded for register the metric only
	metricName := "feed_news_by_category"
	value := 1.0
	timeSeries := createTimeSeries(metricName, value, labels)

	// Create write request
	writeRequest := &prompb.WriteRequest{
		Timeseries: timeSeries,
	}

	// Marshal and compress the request
	data, err := proto.Marshal(writeRequest)
	if err != nil {
		logger.Errorf("Failed to marshal write request: %v", err)
		return err
	}
	compressed := snappy.Encode(nil, data)

	// Write metrics
	if err := handler.WriteMetrics(compressed); err != nil {
		logger.Errorf("Failed to write metrics: %v", err)
		return err
	}

	logger.Debugf("Metrics written successfully for %s", title)
	return nil
}
