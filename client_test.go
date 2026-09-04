package hypixel

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	rateLimit := &RateLimit{}

	client := NewClient(apiKey, rateLimit)

	if client.apiKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.apiKey)
	}

	if client.baseURL != "https://api.hypixel.net/v2/" {
		t.Errorf("Expected base URL %s, got %s", "https://api.hypixel.net/v2/", client.baseURL)
	}

	if client.httpClient != http.DefaultClient {
		t.Error("Expected default HTTP client")
	}

	if client.rate != rateLimit {
		t.Error("rate limit not set correctly")
	}
}

func TestGetters(t *testing.T) {
	customHTTP := &http.Client{}
	client := &Client{
		baseURL:    "https://custom.url/",
		apiKey:     "custom-key",
		httpClient: customHTTP,
		rate:       &RateLimit{},
	}

	t.Run("GetBaseURL", func(t *testing.T) {
		if got := client.GetBaseURL(); got != "https://custom.url/" {
			t.Errorf("GetBaseURL() = %v, want %v", got, "https://custom.url/")
		}
	})

	t.Run("GetAPIKey", func(t *testing.T) {
		if got := client.GetAPIKey(); got != "custom-key" {
			t.Errorf("GetAPIKey() = %v, want %v", got, "custom-key")
		}
	})

	t.Run("GetHTTPClient", func(t *testing.T) {
		if got := client.GetHTTPClient(); got != customHTTP {
			t.Errorf("GetHTTPClient() = %v, want %v", got, customHTTP)
		}
	})
}

func TestSetters(t *testing.T) {
	client := NewClient("initial-key", &RateLimit{})

	t.Run("SetBaseURL", func(t *testing.T) {
		newURL := "https://new.url/"
		client.SetBaseURL(newURL)
		if client.baseURL != newURL {
			t.Errorf("SetBaseURL() failed, got %v, want %v", client.baseURL, newURL)
		}
	})

	t.Run("SetAPIKey", func(t *testing.T) {
		newKey := "new-key"
		client.SetAPIKey(newKey)
		if client.apiKey != newKey {
			t.Errorf("SetAPIKey() failed, got %v, want %v", client.apiKey, newKey)
		}
	})

	t.Run("SetHTTPClient", func(t *testing.T) {
		newClient := &http.Client{}
		client.SetHTTPClient(newClient)
		if client.httpClient != newClient {
			t.Errorf("SetHTTPClient() failed, got %v, want %v", client.httpClient, newClient)
		}
	})
}

func TestClient_Get_URLBuilding(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		expected string
	}{
		{
			name:     "baseURL with trailing slash",
			baseURL:  "https://api.hypixel.net/v2/",
			path:     "skyblock/auctions",
			expected: "https://api.hypixel.net/v2/skyblock/auctions",
		},
		{
			name:     "baseURL without trailing slash",
			baseURL:  "https://api.hypixel.net/v2",
			path:     "skyblock/auctions",
			expected: "https://api.hypixel.net/v2/skyblock/auctions",
		},
		{
			name:     "empty path",
			baseURL:  "https://api.hypixel.net/v2/",
			path:     "",
			expected: "https://api.hypixel.net/v2/",
		},
		{
			name:     "path with leading slash",
			baseURL:  "https://api.hypixel.net/v2/",
			path:     "/skyblock/auctions",
			expected: "https://api.hypixel.net/v2/skyblock/auctions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("test-key", nil)
			client.SetBaseURL(tt.baseURL)
			var gotURL string
			client.SetPreRequestHook(func(r Request) (Response, error) {
				gotURL = r.URL
				return Response{Status: http.StatusOK}, nil
			})
			if _, err := client.Get(Request{Path: tt.path}); err != nil {
				t.Fatalf("Get() returned unexpected error: %v", err)
			}
			if gotURL != tt.expected {
				t.Errorf("Get() built URL = %v, want %v", gotURL, tt.expected)
			}
		})
	}
}
