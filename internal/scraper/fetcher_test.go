package scraper

import (
	"net/url"
	"testing"
	"time"

	"github.com/fromsko/krio/internal/config"
)

func TestNewFetcher(t *testing.T) {
	cfg := &config.ScraperConfig{
		UserAgent:  "test-agent",
		Timeout:     5 * time.Second,
		MaxRetries:  3,
		RetryDelay:  100 * time.Millisecond,
	}

	fetcher := NewFetcher(cfg)
	if fetcher == nil {
		t.Fatal("NewFetcher returned nil")
	}

	if fetcher.cfg != cfg {
		t.Error("fetcher.cfg not set correctly")
	}
}

func TestValidateURL(t *testing.T) {
	cfg := &config.ScraperConfig{
		UserAgent:  "test-agent",
		Timeout:     5 * time.Second,
		MaxRetries:  3,
		RetryDelay:  100 * time.Millisecond,
	}

	fetcher := NewFetcher(cfg)

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid HTTP URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "localhost - should fail",
			url:     "http://localhost:8080",
			wantErr: true,
		},
		{
			name:    "127.0.0.1 - should fail",
			url:     "http://127.0.0.1:8080",
			wantErr: true,
		},
		{
			name:    "private IP 10.x - should fail",
			url:     "http://10.0.0.1:8080",
			wantErr: true,
		},
		{
			name:    "private IP 192.168.x - should fail",
			url:     "http://192.168.1.1:8080",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			url:     "://invalid-url",
			wantErr: true,
		},
		{
			name:    "no protocol",
			url:     "example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fetcher.validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsPrivateURL(t *testing.T) {
	cfg := &config.ScraperConfig{
		UserAgent:  "test-agent",
		Timeout:     5 * time.Second,
		MaxRetries:  3,
		RetryDelay:  100 * time.Millisecond,
	}

	fetcher := NewFetcher(cfg)

	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "localhost",
			url:  "http://localhost",
			want: true,
		},
		{
			name: "127.0.0.1",
			url:  "http://127.0.0.1",
			want: true,
		},
		{
			name: "10.0.0.1",
			url:  "http://10.0.0.1",
			want: true,
		},
		{
			name: "192.168.1.1",
			url:  "http://192.168.1.1",
			want: true,
		},
		{
			name: "172.16.0.1",
			url:  "http://172.16.0.1",
			want: true,
		},
		{
			name: "public IP",
			url:  "http://8.8.8.8",
			want: false,
		},
		{
			name: "domain",
			url:  "http://example.com",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedURL, err := url.Parse(tt.url)
			if err != nil {
				t.Fatalf("Failed to parse URL: %v", err)
			}
			got := fetcher.isPrivateURL(parsedURL)
			if got != tt.want {
				t.Errorf("isPrivateURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
