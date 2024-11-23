package statuspal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// MetricType defines the possible types for metrics (Uptime or Response Time).
type MetricType string

const (
	MetricTypeUptime       MetricType = "up"
	MetricTypeResponseTime MetricType = "rt"
)

// FeaturedNumber defines the possible types for the featured number.
type FeaturedNumber string

const (
	FeaturedNumberAvg  FeaturedNumber = "avg"
	FeaturedNumberMax  FeaturedNumber = "max"
	FeaturedNumberLast FeaturedNumber = "last"
)

// Metric represents a metric on the status page.
type Metric struct {
	ID              int64          `json:"id,omitempty"`
	Status          string         `json:"status,omitempty"`
	LatestEntryTime int64          `json:"latest_entry_time,omitempty"`
	Order           int64          `json:"order,omitempty"`
	Title           string         `json:"title"`
	Unit            string         `json:"unit"`
	Type            MetricType     `json:"type"`
	Enabled         bool           `json:"enabled,omitempty"`
	Visible         bool           `json:"visible,omitempty"`
	RemoteID        string         `json:"remote_id,omitempty"`
	RemoteName      string         `json:"remote_name,omitempty"`
	Threshold       int64          `json:"threshold,omitempty"`
	FeaturedNumber  FeaturedNumber `json:"featured_number,omitempty"`
	IntegrationID   string         `json:"integration_id,omitempty"`
}

type MetricBody struct {
	Metric Metric `json:"metric"`
}

type MetricsBody struct {
	Metrics []Metric `json:"metrics"`
}

type MetricsQuery struct {
	Before string `query:"before"`
	After  string `query:"after"`
	Limit  int64  `query:"limit"`
}

func (c *Client) GetMetrics(statusPageSubdomain string, query MetricsQuery) (*[]Metric, error) {
	urlParams := url.Values{}
	if query.Before != "" {
		urlParams.Add("before", query.Before)
	}
	if query.After != "" {
		urlParams.Add("after", query.After)
	}
	if query.Limit > 0 {
		urlParams.Add("limit", fmt.Sprintf("%d", query.Limit))
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status_pages/%s/metrics%s", c.HostURL, statusPageSubdomain, urlParams.Encode()), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := MetricsBody{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Metrics, nil
}

// GetMetric retrieves a single metric by ID.
func (c *Client) GetMetric(id int64, subdomain string) (*Metric, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status_pages/%s/metrics/%d", c.HostURL, subdomain, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response MetricBody
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Metric, nil
}

// CreateMetric creates a new metric for the status page.
func (c *Client) CreateMetric(subdomain string, metric *Metric) (*Metric, error) {
	rb, err := json.Marshal(MetricBody{
		Metric: *metric,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/status_pages/%s/metrics", c.HostURL, subdomain), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response MetricBody
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Metric, nil
}

// UpdateMetric updates an existing metric on the status page.
func (c *Client) UpdateMetric(id int64, subdomain string, metric *Metric) (*Metric, error) {
	rb, err := json.Marshal(MetricBody{
		Metric: *metric,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/status_pages/%s/metrics/%d", c.HostURL, subdomain, id), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response MetricBody
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Metric, nil
}

// DeleteMetric deletes a metric from the status page.
func (c *Client) DeleteMetric(id int64, subdomain string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/status_pages/%s/metrics/%d", c.HostURL, subdomain, id), nil)
	if err != nil {
		return err
	}

	if _, err := c.doRequest(req); err != nil {
		return err
	}

	return nil
}
