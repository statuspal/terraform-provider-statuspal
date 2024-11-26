package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetricResource(t *testing.T) {
	var integrationID int64 = 1

	create := func(w http.ResponseWriter, r *http.Request) {
		var body statuspal.MetricBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode JSON: %v", err), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		body.Metric.ID = 1
		body.Metric.Status = "active"
		body.Metric.LatestEntryTime = 1633000000
		body.Metric.Order = 1
		body.Metric.Threshold = 100
		body.Metric.FeaturedNumber = "avg"
		body.Metric.IntegrationID = &integrationID

		if err := json.NewEncoder(w).Encode(body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	}

	read := func(w http.ResponseWriter, r *http.Request) {
		metric := statuspal.Metric{
			ID:             1,
			Status:         "active",
			Title:          "Website Response Time",
			Unit:           "ms",
			Type:           "rt",
			Threshold:      100,
			FeaturedNumber: "avg",
			IntegrationID:  &integrationID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		body := statuspal.MetricBody{
			Metric: metric,
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	}

	update := func(w http.ResponseWriter, r *http.Request) {
		metric := statuspal.Metric{
			ID:             1,
			Status:         "active",
			Title:          "Website Response Time",
			Unit:           "ms",
			Type:           "rt",
			Enabled:        true,
			Visible:        true,
			RemoteID:       "9dd9dee9-bf3d-4f2a-9cb6-c8e623b7df5a",
			RemoteName:     "StatusPal",
			Threshold:      100,
			FeaturedNumber: "avg",
			IntegrationID:  nil,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		body := statuspal.MetricBody{
			Metric: metric,
		}

		if err := json.NewEncoder(w).Encode(body); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	}

	remove := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}

	mux := http.NewServeMux()
	mux.Handle("POST /status_pages/{subdomain}/metrics", http.HandlerFunc(create))
	mux.Handle("GET /status_pages/{subdomain}/metrics/{id}", http.HandlerFunc(read))
	mux.Handle("PUT /status_pages/{subdomain}/metrics/{id}", http.HandlerFunc(update))
	mux.Handle("DELETE /status_pages/{subdomain}/metrics/{id}", http.HandlerFunc(remove))

	mock := httptest.NewServer(mux)
	defer mock.Close()

	providerConfig := *providerConfig(&mock.URL)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "statuspal_metric" "test" {
  status_page_subdomain = "example-com-24"
  metric = {
    title             = "Website Response Time"
    unit              = "ms"
    type              = "rt"
  }
}
        `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("statuspal_metric.test", "status_page_subdomain", "example-com-24"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.title", "Website Response Time"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.unit", "ms"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.type", "rt"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.enabled", "false"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.visible", "false"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.remote_id", ""),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.remote_name", ""),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.status", "active"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.latest_entry_time", "1633000000"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.threshold", "100"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.featured_number", "avg"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.order", "1"),
					resource.TestCheckResourceAttr("statuspal_metric.test", "metric.integration_id", "1"),
				),
			},
		},
	})
}
