package statuspal

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client struct.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	ApiKey     string
}

// NewClient function.
func NewClient(api_key *string, region *string, test_url *string) (*Client, error) {
	env := os.Getenv("TF_ENV")

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default StatusPal API URL
		HostURL: "http://local.statuspal.io:4000/api/v2",
	}

	if *region == "EU" || *region == "US" {
		topLevelDomain := map[string]string{
			"EU": "eu",
			"US": "io",
		}[*region]

		c.HostURL = fmt.Sprintf("https://statuspal.%s/api/v2", topLevelDomain)
	} else if env == "TEST" {
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
		return nil, NewError(res.StatusCode, body)
	}

	return &body, nil
}
