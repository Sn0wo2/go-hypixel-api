package hypixel

import (
	"net/http"
	"testing"
)

func TestClient_Authentication(t *testing.T) {
	h := http.Header{}
	h.Set("head", "value1")
	c := NewClient("test1", nil)
	if c.authHeader(h).Get("API-Key") != "test1" {
		t.Errorf("expected 'test1', got %s", h.Get("API-Key"))
	}
	if h.Get("head") != "value1" {
		t.Errorf("expected 'value1', got %s", h.Get("head"))
	}
}

type sampleData struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
	Count   int    `json:"count"`
}

func TestResponse_Decode(t *testing.T) {
	resp := Response{
		Status:  http.StatusOK,
		Content: []byte(`{"success":true,"name":"Hypixel","count":42}`),
	}

	data, err := resp.Decode[sampleData]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !data.Success || data.Name != "Hypixel" || data.Count != 42 {
		t.Errorf("unexpected decoded data: %+v", data)
	}

	badResp := Response{
		Status:  http.StatusOK,
		Content: []byte(`{invalid json`),
	}
	if _, err := badResp.Decode[sampleData](); err == nil {
		t.Error("expected error for invalid json, got nil")
	}
}
