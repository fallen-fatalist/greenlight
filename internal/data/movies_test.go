package data

import (
	"testing"
	"time"

	"greenlight.fallen-fatalist.net/internal/assert"
)

func TestMoviesInsert(t *testing.T) {

	type insertWant struct {
		id      int64
		version int32
	}

	tests := []struct {
		name string
		give Movie
		want insertWant
	}{
		{
			name: "Standard movie",
			give: Movie{
				Title:   "Spider man 2",
				Year:    2004,
				Runtime: 127,
				Genres:  []string{"sci-fi", "action", "adventure"},
			},
			want: insertWant{
				id:      4,
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
			want: insertWant{
				id:      4,
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
			want: insertWant{
				id:      4,
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

func TestMoviesGet(t *testing.T) {
	type getWant struct {
		movie Movie
		err   error
	}
	tests := []struct {
		name string
		give int64
		want getWant
	}{
		{
			name: "1 number ID",
			give: 1,
			want: getWant{
				movie: Movie{
					Title:   "Spider man",
					Year:    2002,
					Runtime: 102,
					Genres:  []string{"sci-fi", "action", "adventure"},
				},
				err: nil,
			},
		},
		{
			name: "2 number ID",
			give: 2,
			want: getWant{
				movie: Movie{
					Title:   "Attack of the titans",
					Year:    2013,
					Runtime: 1,
					Genres:  []string{"adventure", "action", "fantasy", "drama", "cartoon"},
				},
				err: nil,
			},
		},
		{
			name: "3 number ID",
			give: 3,
			want: getWant{
				movie: Movie{
					Title:   "Grimgar of the fantasy and ash",
					Year:    2015,
					Runtime: 1,
					Genres:  []string{"anime", "adventure", "action", "fantasy", "cartoon"},
				},
				err: nil,
			},
		},
		{
			name: "0 number ID",
			give: 0,
			want: getWant{
				movie: Movie{},
				err:   ErrRecordNotFound,
			},
		},
		{
			name: "Non-existing ID",
			give: 5,
			want: getWant{
				movie: Movie{},
				err:   ErrRecordNotFound,
			},
		},
		{
			name: "Negative ID",
			give: -5,
			want: getWant{
				movie: Movie{},
				err:   ErrRecordNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := MovieModel{db}

			movie, err := m.Get(tt.give)
			if movie != nil {
				if !movie.Equal(tt.want.movie) {
					t.Errorf("got: %+v; want %+v", movie, tt.want.movie)
				}
			}
			if err != tt.want.err {
				t.Errorf("got: %+v; want %+v", movie, tt.want.movie)
			}

		})

	}
}

func TestMoviesUpdate(t *testing.T) {
	tests := []struct {
		name string
		give Movie
	}{
		{
			name: "First movie",
			give: Movie{
				ID:      1,
				Title:   "Spider man",
				Year:    2002,
				Runtime: 102,
				Genres:  []string{"sci-fi", "action", "adventure"},
				Version: 1,
			},
		},
		{
			name: "Second movie",
			give: Movie{
				ID:      2,
				Title:   "Attack of the titans",
				Year:    2013,
				Runtime: 1,
				Genres:  []string{"adventure", "action", "fantasy", "drama", "cartoon"},
				Version: 1,
			},
		},
		{
			name: "Third movie",
			give: Movie{
				ID:      3,
				Title:   "Grimgar of the fantasy and ash",
				Year:    2016,
				Runtime: 1,
				Genres:  []string{"anime", "adventure", "action", "fantasy", "cartoon"},
				Version: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := MovieModel{db}

			err := m.Update(&tt.give)
			assert.NilError(t, err)

			if tt.give.Version != 2 {
				t.Errorf("Version of movie after update didn't changed.")
			}

		})

	}

}

func TestMoviesDelete(t *testing.T) {
	tests := []struct {
		name string
		give int64
		want error
	}{
		{
			name: "0 ID",
			give: 0,
			want: ErrRecordNotFound,
		},
		{
			name: "Negative ID",
			give: -1,
			want: ErrRecordNotFound,
		},
		{
			name: "First ID",
			give: 1,
			want: nil,
		},
		{
			name: "Second ID",
			give: 2,
			want: nil,
		},
		{
			name: "Third ID",
			give: 3,
			want: nil,
		},
		{
			name: "Forth ID",
			give: 4,
			want: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := MovieModel{db}

			err := m.Delete(tt.give)
			if err != tt.want {
				t.Errorf("got: %v; want %v", err, tt.want)
			}

		})

	}

}
