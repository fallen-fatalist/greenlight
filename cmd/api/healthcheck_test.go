package main

import (
	"net/http"
	"strings"
	"testing"

	"greenlight.fallen-fatalist.net/internal/assert"
)

func TestHealthCheck(t *testing.T) {
	app := newTestApplication()
	ts := NewTestServer(app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/v1/healthcheck")
	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	expected := `{
		"status": "available",
		"system_info": {
			"environment": "testing",
			"version": "1.0.0"
		}
	}
	`
	expected = strings.ReplaceAll(expected, " ", "")
	expected = strings.ReplaceAll(expected, "\t", "")
	expected = strings.ReplaceAll(expected, "\n", "")

	got := string(body)
	got = strings.ReplaceAll(got, " ", "")
	got = strings.ReplaceAll(got, "\t", "")
	got = strings.ReplaceAll(got, "\n", "")

	assert.Equal(t, got, expected)
}
