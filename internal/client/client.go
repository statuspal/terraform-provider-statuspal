package statuspal

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// HostURL - Default StatusPal API URL
const HostURL string = "http://local.statuspal.io:4000/api/v2"

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

// NewClient -
func NewClient(api_key *string, region *string, test_url *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default StatusPal API URL
		HostURL: HostURL,
	}

	if *region == "eu" || *region == "us" {
		c.HostURL = fmt.Sprintf("https://statuspal.%s/api/v2", *region)
	} else if *region == "test" {
		c.HostURL = *test_url
	}

	// If api_key is not provided, return empty client
	if api_key == nil {
		return &c, nil
	}

	c.ApiKey = *api_key

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) (*[]byte, error) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.ApiKey)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode > http.StatusIMUsed {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return &body, nil
}
