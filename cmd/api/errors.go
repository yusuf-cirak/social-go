package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("Internal server error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Bad request", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusBadRequest, "Bad Request")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Not found", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusNotFound, "Not Found")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "Rate limit exceeded")
}
