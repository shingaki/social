package main

import (
	"net/http"
)

// healthCheck godoc
//
//	@Summary		checks health of the application
//	@Description	checks if API is up and running
//	@Tags			ops
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		APIKeyAuth
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "OK",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
