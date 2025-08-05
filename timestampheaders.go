// Package traefiktimestamping implements a Traefik plugin that adds timestamp headers to HTTP responses.
package traefiktimestamping

import (
	"context"
	"net/http"
	"time"
)

// Config holds the plugin configuration.
type Config struct {
	RequestHeaderName  string `json:"requestHeaderName,omitempty"`
	ResponseHeaderName string `json:"responseHeaderName,omitempty"`
	DateFormat         string `json:"dateFormat,omitempty"`
}

// CreateConfig initializes the plugin configuration with default values.
func CreateConfig() *Config {
	return &Config{
		RequestHeaderName:  "REQUEST-TIMESTAMP",
		ResponseHeaderName: "RESPONSE-TIMESTAMP",
		DateFormat:         "2006-01-02T15:04:05.000Z",
	}
}

// TimestampHeaders implements the Traefik plugin interface.
type TimestampHeaders struct {
	next   http.Handler
	config *Config
}

// New creates a new TimestampHeaders plugin instance.
func New(_ context.Context, next http.Handler, config *Config, _ string) (http.Handler, error) {
	return &TimestampHeaders{
		next:   next,
		config: config,
	}, nil
}

// responseWriter wraps http.ResponseWriter to add both timestamps
type responseWriter struct {
	http.ResponseWriter
	requestTimestamp string
	headerWritten    bool
	config           *Config
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.headerWritten {
		responseTimestamp := time.Now().UTC().Format(rw.config.DateFormat)
		rw.Header().Set(rw.config.RequestHeaderName, rw.requestTimestamp)
		rw.Header().Set(rw.config.ResponseHeaderName, responseTimestamp)
		rw.headerWritten = true
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (t *TimestampHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	requestTimestamp := time.Now().UTC().Format(t.config.DateFormat)

	// Wrap the response writer to add both timestamps when response is sent
	wrappedWriter := &responseWriter{
		ResponseWriter:   rw,
		requestTimestamp: requestTimestamp,
		headerWritten:    false,
		config:           t.config,
	}

	t.next.ServeHTTP(wrappedWriter, req)
}
