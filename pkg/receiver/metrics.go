package receiver

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

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
	req, err := http.NewRequest("POST", r.url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	time.Sleep(1 * time.Second)
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
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

func produceMetricsToPrometheus(username string, password string, prometheusURL string, category string, title string, url string, pubDate string) {
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
		log.Fatalf("failed to marshal write request: %v", err)
	}
	compressed := snappy.Encode(nil, data)

	// Write metrics
	err = handler.WriteMetrics(compressed)
	if err != nil {
		log.Fatalf("failed to write metrics with status: %v", err)
	}

	log.Println("Metrics written to remote Prometheus server successfully")
}
