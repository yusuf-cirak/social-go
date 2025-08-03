package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	writeJSON(w, http.StatusOK, data)

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
	}
}
