package main

import (
	"net/http"

	"github.com/yusuf-cirak/social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	posts, err := app.store.Posts.GetUserFeed(ctx, int64(42), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
