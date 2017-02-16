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
