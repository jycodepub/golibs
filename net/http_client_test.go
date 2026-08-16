package net

import (
	"testing"
)

func TestHttpClient_GetGoogle(t *testing.T) {
	client := NewHttpClient()
	resp, err := client.Get("http://www.google.com", nil)
	if err != nil {
		t.Fatalf("failed to GET http://www.google.com: %v", err)
	}

	if resp.Code != 200 {
		t.Errorf("expected response code 200, got %d", resp.Code)
	}

	if !resp.IsOK() {
		t.Errorf("expected resp.IsOK() to be true, got false")
	}

	if len(resp.Body) == 0 {
		t.Errorf("expected non-empty response body")
	}
}
