package main

import (
	"net/http"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	status := envelope{
		"system_info": map[string]string{
			"Status":      "Active",
			"Version":     Version,
			"Environment": app.config.env,
		},
	}

	err := app.writeJson(w, http.StatusOK, status, nil)
	if err != nil {
	    app.serverErrorResponse(w, r, err)	
	}
}
