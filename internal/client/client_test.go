package statuspal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_doRequest_rate_limit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{
		HostURL:    server.URL,
		HTTPClient: server.Client(),
		ApiKey:     "test",
	}

	start := time.Now()
	var requestsCount uint

	// NOTE: Executes the request 100 times to have enough register to validate the request per second.
	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest(http.MethodGet, client.HostURL, nil)
		if _, err := client.doRequest(req.WithContext(context.Background())); err != nil {
			t.Errorf("Request error at iteration %d: %v", i, err)
		}

		requestsCount += 1
	}

	requestsTime := time.Since(start).Round(time.Duration(RateLimiter.Limit()) * time.Second)

	requestPerSecond := requestsCount / uint(requestsTime.Seconds())
	if requestPerSecond != uint(RateLimiter.Burst()) {
		t.Fatal("Requests rate doesn't meet the request rate required")
	}

	t.Logf("All requests executed within the rate limit of %d per second", requestPerSecond)
}
