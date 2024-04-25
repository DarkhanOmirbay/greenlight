package main

import (
	"errors"
	"fmt"
	"greenlight.darkhanomirbay/internal/data"
	"greenlight.darkhanomirbay/internal/validator"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	// json.Unmarshal!!!
	//body, err := io.ReadAll(r.Body)
	//if err != nil {
	//	app.serverErrorResponse(w, r, err)
	//	return
	//}
	//err = json.Unmarshal(body, &input)
	//if err != nil {
	//	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	//}
	//err := json.NewDecoder(r.Body).Decode(&input)
	err := app.readJSON(w, r, &input)
	if err != nil {
		//app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		app.badRequestResponse(w, r, err)
	}
	v := validator.New()
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres}
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	//fmt.Fprintf(w, "%+v\n", input)
	//v.Check(input.Title != "", "title", "must be provided")
	//v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")
	//
	//v.Check(input.Year != 0, "year", "must be provided")
	//v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	//v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future ")
	//
	//v.Check(input.Runtime != 0, "runtime", "must be provided")
	//v.Check(input.Runtime > 0, "runtime", "must be a postive integer")
	//
	//v.Check(input.Genres != nil, "genres", "must be provided")
	//v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	//v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genre")
	//v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")
	//if !v.Valid() {
	//	app.failedValidationResponse(w, r, v.Errors)
	//	return
	//}

}
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//movie := data.Movie{
	//	ID:        id,
	//	CreatedAt: time.Now(),
	//	Title:     "Casablanca",
	//	Runtime:   102,
	//	Genres:    []string{"drama", "romance", "war"},
	//	Version:   1,
	//}
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}
	//fmt.Fprintf(w, "movie id is %d", id)
}
func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}
	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}
	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return

	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string
		Genres  []string
		Filters data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCsv(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	fmt.Fprintf(w, "%+v\n", input)
}
