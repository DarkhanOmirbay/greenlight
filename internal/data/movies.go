package data

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"greenlight.darkhanomirbay/internal/validator"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // - (hyphen)directive use for hiding field
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty,string"` // string directive use for represent field in JSON STRING
	// key-word omitempty uses for hide empty field
	Genres  []string `json:"genres,omitempty"`
	Version int32    `json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

func (m *MovieModel) Insert(movie *Movie) error {
	query := `INSERT INTO movies(title,year,runtime,genres) VALUES($1,$2,$3,$4) RETURNING id,created_at,version `
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}
func (m *MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, created_at, title, year, runtime, genres, version
FROM movies
WHERE id = $1`

	var movie Movie
	err := m.DB.QueryRow(query, id).Scan( // Use scan for save fields into movie(copy)
		&movie.ID, // use & this for record field into movie
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres), // pq.Array for convert text[] PostgresSQL's field to our array
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}
	return &movie, err
}
func (m *MovieModel) Update(movie *Movie) error {
	query := `UPDATE movies SET title=$1,year=$2,runtime=$3,genres=$4,version=version+1 WHERE id=$5 RETURNING version `
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID}
	return m.DB.QueryRow(query, args...).Scan(&movie.Version)

}
func (m *MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM movies WHERE id=$1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

// MOCK MODELS (FOR UNIT TESTS)
//type MockMovieModel struct{}
//
//func (m MockMovieModel) Insert(movie *Movie) error {
//	// Mock the action...
//	return nil
//}
//func (m MockMovieModel) Get(id int64) (*Movie, error) {
//	// Mock the action...
//	return nil, nil
//}
//func (m MockMovieModel) Update(movie *Movie) error {
//	// Mock the action...
//	return nil
//}
//func (m MockMovieModel) Delete(id int64) error {
//	// Mock the action...
//	return nil
//}

//type Movie struct {
//	ID        int64
//	CreatedAt time.Time
//	Title     string
//	Year      int32
//	Runtime   int32
//	Genres    []string
//	Version   int32
//}

//func (m Movie) MarshalJSON() ([]byte, error) {
//
//	var runtime string
//
//	if m.Runtime != 0 {
//		runtime = fmt.Sprintf("%d mins", m.Runtime)
//	}
//	aux := struct {
//		ID      int64    `json:"id"`
//		Title   string   `json:"title"`
//		Year    int32    `json:"year,omitempty"`
//		Runtime string   `json:"runtime,omitempty"` // This is a string.
//		Genres  []string `json:"genres,omitempty"`
//		Version int32    `json:"version"`
//	}{
//		ID:      m.ID,
//		Title:   m.Title,
//		Year:    m.Year,
//		Runtime: runtime, // Note that we assign the value from the runtime variable here.
//		Genres:  m.Genres,
//		Version: m.Version,
//	}
//	return json.Marshal(aux)
//}

//func (m Movie) MarshalJSON() ([]byte, error) {
//	// Create a variable holding the custom runtime string, just like before.
//	var runtime string
//	if m.Runtime != 0 {
//		runtime = fmt.Sprintf("%d mins", m.Runtime)
//	}
//	// Define a MovieAlias type which has the underlying type Movie. Due to the way that
//	// Go handles type definitions (https://golang.org/ref/spec#Type_definitions) the
//	// MovieAlias type will contain all the fields that our Movie struct has but,
//	// importantly, none of the methods.
//	type MovieAlias Movie
//	// Embed the MovieAlias type inside the anonymous struct, along with a Runtime field
//	// that has the type string and the necessary struct tags. It's important that we
//	// embed the MovieAlias type here, rather than the Movie type directly, to avoid
//	// inheriting the MarshalJSON() method of the Movie type (which would result in an
//	// infinite loop during encoding).
//	aux := struct {
//		MovieAlias
//		Runtime string `json:"runtime,omitempty"`
//	}{
//		MovieAlias: MovieAlias(m),
//		Runtime: runtime,
//	}
//	return json.Marshal(aux)
//}
