package collector

import (
	"encoding/json"
	"net"
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
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(network, addr, timeout)
					if err != nil {
						return nil, err
					}
					if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
						return nil, err
					}
					return c, nil
				},
			},
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
