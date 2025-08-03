package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal error %s path: %s error: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request %s path: %s error: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, "Bad Request")
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found %s path: %s error: %v", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "Not Found")
}
