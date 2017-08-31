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
	resp, err := c.httpClient.Get(c.nsqdURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sr statsResponse
	if err = json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, err
	}
	return &sr.Data, nil
}
