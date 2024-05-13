package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"greenlight.fallen-fatalist.net/internal/data"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication() *application {
	cfg := config{env: "testing"}

	return &application{
		cfg:    cfg,
		models: data.NewMockModels(),
	}
}

func NewTestServer(h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := rs.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body

}
