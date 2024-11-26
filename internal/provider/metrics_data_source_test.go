package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	statuspal "terraform-provider-statuspal/internal/client"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetricsDataSource(t *testing.T) {
	integrationID := int64(1)

	listSingleMetric := func(w http.ResponseWriter, r *http.Request) {
		metric := statuspal.Metric{
			ID:              1,
			Status:          "active",
			Title:           "Website Response Time",
			Unit:            "ms",
			Type:            "rt",
			Threshold:       100,
			FeaturedNumber:  "avg",
			IntegrationID:   &integrationID,
			LatestEntryTime: 1633000000,
			Order:           1,
		}
		metrics := []statuspal.Metric{metric}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(statuspal.MetricsBody{Metrics: metrics}) //nolint:errcheck
	}

	listMultipleMetrics := func(w http.ResponseWriter, r *http.Request) {
		metrics := []statuspal.Metric{
			{
				ID:              1,
				Status:          "active",
				Title:           "Website Response Time",
				Unit:            "ms",
				Type:            "rt",
				Threshold:       100,
				FeaturedNumber:  "avg",
				IntegrationID:   &integrationID,
				LatestEntryTime: 1633000000,
				Order:           1,
			},
			{
				ID:              2,
				Status:          "inactive",
				Title:           "CPU Usage",
				Unit:            "%",
				Type:            "gauge",
				Threshold:       90,
				FeaturedNumber:  "max",
				IntegrationID:   nil,
				LatestEntryTime: 1633000100,
				Order:           2,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(statuspal.MetricsBody{Metrics: metrics}) //nolint:errcheck
	}

	listEmptyMetrics := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(statuspal.MetricsBody{Metrics: []statuspal.Metric{}}) //nolint:errcheck
	}

	listError := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{ //nolint:errcheck
			"error": "Internal Server Error",
		})
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/status_pages/example-com-24/metrics", listSingleMetric)
	mux.HandleFunc("/status_pages/example-com-25/metrics", listMultipleMetrics)
	mux.HandleFunc("/status_pages/example-com-26/metrics", listEmptyMetrics)
	mux.HandleFunc("/status_pages/example-com-27/metrics", listError)

	mock := httptest.NewServer(mux)
	defer mock.Close()
	providerConfig := *providerConfig(&mock.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
data "statuspal_metrics" "test" {
  status_page_subdomain = "example-com-24"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "status_page_subdomain", "example-com-24"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.title", "Website Response Time"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.unit", "ms"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.type", "rt"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.status", "active"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.latest_entry_time", "1633000000"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.threshold", "100"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.featured_number", "avg"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.order", "1"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test", "metrics.0.integration_id", "1"),
				),
			},
			{
				Config: providerConfig + `
data "statuspal_metrics" "test_multiple" {
  status_page_subdomain = "example-com-25"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_multiple", "status_page_subdomain", "example-com-25"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_multiple", "metrics.#", "2"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_multiple", "metrics.0.title", "Website Response Time"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_multiple", "metrics.1.title", "CPU Usage"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_multiple", "metrics.1.unit", "%"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_multiple", "metrics.1.type", "gauge"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_multiple", "metrics.1.integration_id", "0"),
				),
			},
			{
				Config: providerConfig + `
data "statuspal_metrics" "test_empty" {
  status_page_subdomain = "example-com-26"
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_empty", "status_page_subdomain", "example-com-26"),
					resource.TestCheckResourceAttr("data.statuspal_metrics.test_empty", "metrics.#", "0"),
				),
			},
		},
	})
}
