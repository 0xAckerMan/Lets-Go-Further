package main

import "net/http"

func (app *Application) logError (r *http.Request, err error){
    app.logger.Println(err)
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}){
    env := envelope{"error": message}
}
