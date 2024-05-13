package main

import (
	"net/http"
	"testing"

	"greenlight.fallen-fatalist.net/internal/assert"
)

func TestShowMovieHandler(t *testing.T) {
	app := newTestApplication()

	ts := NewTestServer(app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/v1/movies/1",
			wantCode: http.StatusOK,
			wantBody: "Spider man",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/v1/movies/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, string(body), tt.wantBody)
			}
		})
	}
}
