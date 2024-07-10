package statuspal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type servicesResponse struct {
	Services []Service `json:"services"`
}

type ServiceResponse struct {
	Service Service `json:"service"`
}

// GetService - Returns list of services from the status page.
func (c *Client) GetServices(statusPageSubdomain *string) (*[]Service, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/status_pages/%s/services", c.HostURL, *statusPageSubdomain), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	response := servicesResponse{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Services, nil
}

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

	response := ServiceResponse{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Service, nil
}

// CreateService - Create new service in the organization.
func (c *Client) CreateService(service *Service, statusPageSubdomain *string) (*Service, error) {
	rb, err := json.Marshal(ServiceResponse{Service: *service})
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

	response := ServiceResponse{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Service, nil
}

// UpdateService - Update a service in the organization.
func (c *Client) UpdateService(service *Service, statusPageSubdomain *string, serviceID *string) (*Service, error) {
	rb, err := json.Marshal(ServiceResponse{Service: *service})
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

	response := ServiceResponse{}
	err = json.Unmarshal(*body, &response)
	if err != nil {
		return nil, err
	}

	return &response.Service, nil
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
