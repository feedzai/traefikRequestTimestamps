package traefikRequestTimestamps

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimestampHeaders(t *testing.T) {
	cfg := CreateConfig()
	ctx := context.Background()

	// Mock backend handler that writes a response
	next := http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("Hello World"))
	})

	plugin, err := New(ctx, next, cfg, "timestamp-headers")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	plugin.ServeHTTP(recorder, req)

	// Validate REQUEST-TIMESTAMP header
	requestTimestamp := recorder.Header().Get("REQUEST-TIMESTAMP")
	if requestTimestamp == "" {
		t.Fatalf("Expected 'REQUEST-TIMESTAMP' header in response, but got none")
	}

	// Validate RESPONSE-TIMESTAMP header
	responseTimestamp := recorder.Header().Get("RESPONSE-TIMESTAMP")
	if responseTimestamp == "" {
		t.Fatalf("Expected 'RESPONSE-TIMESTAMP' header in response, but got none")
	}

	// Parse and validate the request timestamp format
	_, err = time.Parse("2006-01-02T15:04:05.000Z", requestTimestamp)
	if err != nil {
		t.Errorf("Invalid request timestamp format: %v", err)
	}

	// Parse and validate the response timestamp format
	_, err = time.Parse("2006-01-02T15:04:05.000Z", responseTimestamp)
	if err != nil {
		t.Errorf("Invalid response timestamp format: %v", err)
	}

	// Validate that response timestamp is after or equal to request timestamp
	reqTime, _ := time.Parse("2006-01-02T15:04:05.000Z", requestTimestamp)
	respTime, _ := time.Parse("2006-01-02T15:04:05.000Z", responseTimestamp)

	if respTime.Before(reqTime) {
		t.Errorf("Response timestamp (%s) should be after request timestamp (%s)", responseTimestamp, requestTimestamp)
	}
}

func TestTimestampHeadersWithCustomConfig(t *testing.T) {
	// Custom configuration with different header names and date format
	cfg := &Config{
		RequestHeaderName:  "X-Request-Time",
		ResponseHeaderName: "X-Response-Time",
		DateFormat:         "2006-01-02 15:04:05 UTC",
	}
	ctx := context.Background()

	// Mock backend handler that writes a response
	next := http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
		_, _ = rw.Write([]byte("Hello World"))
	})

	plugin, err := New(ctx, next, cfg, "timestamp-headers")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	plugin.ServeHTTP(recorder, req)

	// Validate custom REQUEST header
	requestTimestamp := recorder.Header().Get("X-Request-Time")
	if requestTimestamp == "" {
		t.Fatalf("Expected 'X-Request-Time' header in response, but got none")
	}

	// Validate custom RESPONSE header
	responseTimestamp := recorder.Header().Get("X-Response-Time")
	if responseTimestamp == "" {
		t.Fatalf("Expected 'X-Response-Time' header in response, but got none")
	}

	// Validate that default headers are NOT present
	if recorder.Header().Get("REQUEST-TIMESTAMP") != "" {
		t.Errorf("Should not have default 'REQUEST-TIMESTAMP' header when custom config is used")
	}
	if recorder.Header().Get("RESPONSE-TIMESTAMP") != "" {
		t.Errorf("Should not have default 'RESPONSE-TIMESTAMP' header when custom config is used")
	}

	// Parse and validate the custom timestamp format
	_, err = time.Parse("2006-01-02 15:04:05 UTC", requestTimestamp)
	if err != nil {
		t.Errorf("Invalid request timestamp format: %v", err)
	}

	_, err = time.Parse("2006-01-02 15:04:05 UTC", responseTimestamp)
	if err != nil {
		t.Errorf("Invalid response timestamp format: %v", err)
	}
}
