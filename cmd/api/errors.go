package main

import (
	"fmt"
	"net/http"
)

func (app *Application) logError (r *http.Request, err error){
    app.logger.Println(err)
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}){
    env := envelope{"error": message}
    err := app.writeJson(w, http.StatusOK, env, nil)
    if err != nil{
        app.logError(r, err)
        w.WriteHeader(500)
    }
}

func (app*Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error){
    app.logError(r,err)
    message :="Sorry the server experienced an error, we cant complete your request at the moment"
    app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request){
    message := "the requested resource could not be found"
    app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request){
    message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)

    app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string){
    app.errorResponse(w,r,http.StatusUnprocessableEntity, errors)
}

