package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/0xAckerMan/Lets-Go-Further/internal/data"
)

func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime data.Runtime    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
        app.badRequestResponse(w,r,err)
        return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := &data.Movie{
		Id:        id,
		CreatedAt: time.Now(),
		Title:     "Click Click Bang",
		Runtime:   180,
		Genres:    []string{"drama", "romance", "crime"},
		Version:   1,
	}

	err = app.writeJson(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
