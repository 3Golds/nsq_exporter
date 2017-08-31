package collector

import (
	"encoding/json"
	"net/http"
	"time"
)

type nsqStatsClient struct {
	nsqdURL    string
	httpClient *http.Client
}

func newNSQStatsClient(nsqdURL string, timeout time.Duration) *nsqStatsClient {
	return &nsqStatsClient{
		nsqdURL: nsqdURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *nsqStatsClient) getStats() (*stats, error) {
	req, err := http.NewRequest("GET", c.nsqdURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.nsq; version=1.0")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var statsData stats
	if err = json.NewDecoder(resp.Body).Decode(&statsData); err != nil {
		return nil, err
	}
	return &statsData, nil
}
