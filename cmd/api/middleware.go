package main

import "net/http"

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
			app.rateLimitExceededResponse(w, r, retryAfter.String())
			return
		}
		next.ServeHTTP(w, r)
	})
}
