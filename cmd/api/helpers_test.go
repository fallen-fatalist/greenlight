package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/julienschmidt/httprouter"
	"greenlight.fallen-fatalist.net/internal/assert"
	"greenlight.fallen-fatalist.net/internal/data"
)

func TestFormatInt64(t *testing.T) {
	tests := []struct {
		name string
		give int64
		want string
	}{
		{
			name: "1 number test",
			give: 1,
			want: "1",
		},
		{
			name: "0 number test",
			give: 0,
			want: "0",
		},
		{
			name: "negative number test",
			give: -100,
			want: "-100",
		},
		{
			name: "1000 number test",
			give: 1000,
			want: "1000",
		},
		{
			name: "Maximum int64 number test",
			give: 1<<63 - 1,
			want: "9223372036854775807",
		},
		{
			name: "Minimum int64 number test",
			give: -(1 << 63),
			want: "-9223372036854775808",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, strconv.FormatInt(tt.give, 10), tt.want)
		})
	}

}

func TestReadIDParam(t *testing.T) {
	app := newTestApplication()

	tests := []struct {
		name   string
		giveID int64
		wantID int64
	}{
		{
			name:   "1 number test",
			giveID: 1,
			wantID: 1,
		},
		{
			name:   "0 number test",
			giveID: 0,
			wantID: 0,
		},
		{
			name:   "negative number test",
			giveID: -100,
			wantID: 0,
		},
		{
			name:   "1000 number test",
			giveID: 1000,
			wantID: 1000,
		},
		{
			name:   "Maximum int64 number test",
			giveID: 1<<63 - 1,
			wantID: 1<<63 - 1,
		},
		{
			name:   "Minimum int64 number test",
			giveID: -(1 << 63),
			wantID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/movies/"+strconv.FormatInt(tt.giveID, 10), nil)
			params := httprouter.Params{{Key: "id", Value: strconv.FormatInt(tt.giveID, 10)}}
			ctx := req.Context()
			ctx = context.WithValue(ctx, httprouter.ParamsKey, params)
			req = req.WithContext(ctx)

			id, _ := app.readIDParam(req)
			assert.Equal(t, id, tt.wantID)
		})
	}

}

func TestReadJSON(t *testing.T) {
	app := newTestApplication()

	tests := []struct {
		name string
		give data.Movie
	}{
		{
			name: "Standard test by movie struct",
			give: data.Movie{
				Title:   "Spidi racer",
				Year:    2000,
				Runtime: 100,
				Genres:  []string{"sci-fi", "race"},
			},
		},
		{
			name: "Only title test by movie struct",
			give: data.Movie{
				Title: "Spider man",
			},
		},
		{
			name: "Only title and year test by movie struct",
			give: data.Movie{
				Title: "Spider man 2",
				Year:  2012,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js, err := json.Marshal(tt.give)
			if err != nil {
				t.Fatalf("failed to marshal JSON: %v", err)
			}

			var input struct {
				ID      int64        `json: "id"`
				Version int32        `json: "version"`
				Title   string       `json: "title"`
				Year    int32        `json: "year"`
				Runtime data.Runtime `json: "runtime"`
				Genres  []string     `json: "genres"`
			}

			writeRecorder := httptest.NewRecorder()
			readRequest := httptest.NewRequest("GET", "/v1/movies/", bytes.NewBuffer(js))

			err = app.readJSON(writeRecorder, readRequest, &input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			readMovie := data.Movie{
				Title:   input.Title,
				Year:    input.Year,
				Runtime: input.Runtime,
				Genres:  input.Genres,
			}

			if !tt.give.Equal(readMovie) {
				t.Errorf("got: %v; want %v", readMovie, tt.give)
			}
		})
	}

}

func TestWriteJSON(t *testing.T) {
	app := newTestApplication()

	tests := []struct {
		name string
		give envelope
	}{
		{
			name: "Standard test by movie struct",
			give: envelope{
				"movie": data.Movie{
					Title:   "Spidi racer",
					Year:    2000,
					Runtime: 100,
					Genres:  []string{"sci-fi", "race"},
				}},
		},
		{
			name: "Only title and year test by movie struct",
			give: envelope{
				"movie": data.Movie{
					Title: "Spider man 2",
					Year:  2012,
				}},
		},
		{
			name: "Only title test by movie struct",
			give: envelope{
				"movie": data.Movie{
					Title: "Spider man",
				}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			if err := app.writeJSON(rr, http.StatusOK, tt.give, nil); err != nil {
				t.Fatalf("writeJSON returned an unexpected error: %v", err)
			}

			if rr.Code != http.StatusOK {
				t.Errorf("unexpected status code: got %d, want %d", rr.Code, http.StatusOK)
			}

			if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("unexpected content type: got %s, want application/json", contentType)
			}

			var responseBody envelope
			if err := json.Unmarshal(rr.Body.Bytes(), &responseBody); err != nil {
				t.Fatalf("failed to unmarshal response JSON: %v", err)
			}

			if _, ok := responseBody["movie"]; !ok {
				t.Fatalf("no movie key in response body")
			}

			var movieMap map[string]interface{}
			var ok bool
			if movieMap, ok = responseBody["movie"].(map[string]interface{}); !ok {
				t.Fatalf("response body movie is not the map")
			}

			var input struct {
				ID      int64        `json: "id",omitempty`
				Version int32        `json: "version",omitempty`
				Title   string       `json: "title"`
				Year    int32        `json: "year"`
				Runtime data.Runtime `json: "runtime"`
				Genres  []string     `json: "genres"`
			}

			js, err := json.Marshal(movieMap)
			if err != nil {
				t.Fatal("error while marshalling movie map")
			}
			json.Unmarshal(js, &input)

			movie := data.Movie{
				Year:    input.Year,
				Title:   input.Title,
				Runtime: input.Runtime,
				Genres:  input.Genres,
			}

			if _, ok := tt.give["movie"]; !ok {
				panic(fmt.Sprintf("no movie key in literal envelope"))
			}

			if !movie.Equal(tt.give["movie"].(data.Movie)) {
				t.Errorf("got: %+v; want: %+v", movie, tt.give["movie"])
			}

		})
	}

}
