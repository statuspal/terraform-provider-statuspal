package statuspal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type statusPagesResponse struct {
	StatusPages []StatusPage `json:"status_pages"`
}

type statusPageResponse struct {
	StatusPage StatusPage `json:"status_page"`
}

// GetStatusPages - Returns list of status pages from the organization
func (c *Client) GetStatusPages(organizationID *string) (*[]StatusPage, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/orgs/%s/status_pages", c.HostURL, *organizationID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := statusPagesResponse{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.StatusPages, nil
}

// GetStatusPage - Returns specific status page from the organization
func (c *Client) GetStatusPage(organizationID *string, statusPageSubdomain *string) (*StatusPage, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/orgs/%s/status_pages/%s", c.HostURL, *organizationID, *statusPageSubdomain), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := statusPageResponse{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.StatusPage, nil
}

// CreateStatusPage - Create new status page in the organization
func (c *Client) CreateStatusPage(statusPage *StatusPage, organizationID *string) (*StatusPage, error) {
	rb, err := json.Marshal(statusPageResponse{StatusPage: *statusPage})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/orgs/%s/status_pages", c.HostURL, *organizationID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := statusPageResponse{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.StatusPage, nil
}

// UpdateStatusPage - Update a status page in the organization
func (c *Client) UpdateStatusPage(statusPage *StatusPage, organizationID *string, statusPageSubdomain *string) (*StatusPage, error) {
	rb, err := json.Marshal(statusPageResponse{StatusPage: *statusPage})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/orgs/%s/status_pages/%s", c.HostURL, *organizationID, *statusPageSubdomain), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := statusPageResponse{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.StatusPage, nil
}

// DeleteStatusPage - Delete a status page in the organization
func (c *Client) DeleteStatusPage(organizationID *string, statusPageSubdomain *string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/orgs/%s/status_pages/%s", c.HostURL, *organizationID, *statusPageSubdomain), nil)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	convertedBody := string(*body)
	if convertedBody != `""` {
		return errors.New(convertedBody)
	}

	return nil
}
