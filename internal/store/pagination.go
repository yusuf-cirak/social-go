package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1, lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	offset := qs.Get("offset")
	sort := qs.Get("sort")

	if limit == "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}

		fq.Limit = l
	}

	if offset == "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}

		fq.Offset = o
	}

	if sort == "" {
		fq.Sort = "asc"
	} else {
		fq.Sort = sort
	}

	return fq, nil

}
