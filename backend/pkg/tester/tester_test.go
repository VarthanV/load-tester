package tester

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/VarthanV/load-tester/pkg/liveupdate"
)

type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.RoundTripFunc == nil {
		return nil, fmt.Errorf("no RoundTripFunc defined")
	}
	return m.RoundTripFunc(req)
}

func TestProcessStat(t *testing.T) {
	driver := &driver{
		responseTimeInSeconds: []float64{},
	}

	stat := &RequestStat{IsSuccess: true, TimeTakenInSeconds: 5}
	driver.processStat(stat)

	if driver.requestsSucceeded.Load() != 1 {
		t.Errorf("expected 1 successful request, got %d", driver.requestsSucceeded.Load())
	}

	if driver.requestsFailed.Load() != 0 {
		t.Errorf("expected 0 failed requests, got %d", driver.requestsFailed.Load())
	}

	if len(driver.responseTimeInSeconds) != 1 || driver.responseTimeInSeconds[0] != 5 {
		t.Errorf("unexpected response times: %v", driver.responseTimeInSeconds)
	}
}

func TestNewDriver(t *testing.T) {
	driver, err := New(
		liveupdate.New(),
		WithPeakConfig(50, 10*time.Minute, 10),
		WithRequestConfig("http://example.com", map[string]string{"key": "value"}, http.StatusOK),
		WithHeaders(map[string]string{"Content-Type": "application/json"}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if driver.TargetUsers != 50 {
		t.Errorf("expected TargetUsers to be 50, got %d", driver.TargetUsers)
	}

	if driver.URL != "http://example.com" {
		t.Errorf("expected URL to be http://example.com, got %s", driver.URL)
	}

	if driver.Headers.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type header to be application/json, got %s", driver.Headers.Get("Content-Type"))
	}
}

func TestDoRequestAndReturnStats(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       http.NoBody,
			}, nil
		},
	}

	mockClient := &http.Client{Transport: mockTransport}

	driver := &driver{
		httpClient: mockClient,
		config: config{
			Method:             "GET",
			URL:                "http://example.com",
			SuccessStatusCodes: []int{http.StatusOK},
		},
	}

	stat, err := driver.doRequestAndReturnStats(context.Background(), "GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !stat.IsSuccess {
		t.Errorf("expected IsSuccess to be true, but got false")
	}
}

func TestDoRequestFailure(t *testing.T) {
	mockTransport := &MockRoundTripper{
		RoundTripFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	mockClient := &http.Client{Transport: mockTransport}

	driver := &driver{
		httpClient: mockClient,
		config: config{
			Method:             "GET",
			URL:                "http://example.com",
			SuccessStatusCodes: []int{http.StatusOK},
		},
	}

	stat, err := driver.doRequestAndReturnStats(context.Background(), "GET", "http://example.com", nil)
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if stat != nil && stat.IsSuccess {
		t.Errorf("expected IsSuccess to be false, but got true")
	}
}
