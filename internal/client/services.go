package statuspal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// GetService - Returns specific service from the organization.
func (c *Client) GetService(statusPageSubdomain *string, serviceID *string) (*Service, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status_pages/%s/services/%s", c.HostURL, *statusPageSubdomain, *serviceID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	service := Service{}
	err = json.Unmarshal(*body, &service)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

// CreateService - Create new service in the organization.
func (c *Client) CreateService(service *Service, statusPageSubdomain *string) (*Service, error) {
	rb, err := json.Marshal(*service)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/status_pages/%s/services", c.HostURL, *statusPageSubdomain), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	newService := Service{}
	err = json.Unmarshal(*body, &newService)
	if err != nil {
		return nil, err
	}

	return &newService, nil
}

// UpdateService - Update a service in the organization.
func (c *Client) UpdateService(service *Service, statusPageSubdomain *string, serviceID *string) (*Service, error) {
	rb, err := json.Marshal(*service)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/status_pages/%s/services/%s", c.HostURL, *statusPageSubdomain, *serviceID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	updatedService := Service{}
	err = json.Unmarshal(*body, &updatedService)
	if err != nil {
		return nil, err
	}

	return &updatedService, nil
}

// DeleteService - Delete a service in the organization.
func (c *Client) DeleteService(statusPageSubdomain *string, serviceID *string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/status_pages/%s/services/%s", c.HostURL, *statusPageSubdomain, *serviceID), nil)
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
