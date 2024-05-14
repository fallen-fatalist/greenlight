package data

import (
	"testing"
	"time"

	"greenlight.fallen-fatalist.net/internal/assert"
)

func TestMoviesInsert(t *testing.T) {

	type want struct {
		id      int64
		version int32
	}

	tests := []struct {
		name string
		give Movie
		want want
	}{
		{
			name: "Standard movie",
			give: Movie{
				Title:   "Spider man 2",
				Year:    2004,
				Runtime: 127,
				Genres:  []string{"sci-fi", "action", "adventure"},
			},
			want: want{
				id:      1,
				version: 1,
			},
		},
		{
			name: "Second standard movie",
			give: Movie{
				Title:   "Source code",
				Year:    2011,
				Runtime: 93,
				Genres:  []string{"sci-fi", "thriller", "drama"},
			},
			want: want{
				id:      1,
				version: 1,
			},
		},
		{
			name: "Movie with only title",
			give: Movie{
				Title: "Unnamed",
			},
		},
		{
			name: "Movie with only title and year",
			give: Movie{
				Title: "Unnamed",
				Year:  2000,
			},
		},
		{
			name: "Movie without title",
			give: Movie{
				Year:    2000,
				Runtime: 111,
				Genres:  []string{"sci-fi", "action", "anime"},
			},
			want: want{
				id:      1,
				version: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := MovieModel{db}

			m.Insert(&tt.give)

			assert.Equal(t, tt.give.ID, tt.want.id)
			assert.Equal(t, tt.give.Version, tt.want.version)
			if !tt.give.CreatedAt.IsZero() {
				if !time.Now().Truncate(time.Minute).Equal(tt.give.CreatedAt.Truncate(time.Minute)) {
					t.Errorf("different time of creation; got: %v, expected: %v",
						time.Now().Truncate(time.Minute),
						tt.give.CreatedAt.Truncate(time.Minute),
					)
				}
			}
		})

	}
}
