package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yusuf-cirak/social/internal/auth"
	"github.com/yusuf-cirak/social/internal/store"
	"go.uber.org/zap"
)

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger // zap.Logger is much faster but only does structured logging.
	jwt    *auth.Manager
	policy *auth.PolicyEngine
}

type config struct {
	addr string
	db   dbConfig
	env  string
	auth authConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {

		r.Get("/health", app.healthCheckHandler)

		// auth endpoints
		r.Post("/auth/login", app.loginHandler)

		r.Route("/posts", func(r chi.Router) {
			// Protected routes with authorization
			r.Group(func(r chi.Router) {
				r.Use(app.authMiddleware)
				r.Use(app.authorize(auth.ActionPostCreate, app.resourcePostCreate))
				r.Post("/", app.createPostHandler)
			})
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)

				r.Get("/", app.getPostHandler)

				// Delete post - requires auth + ownership
				r.Group(func(r chi.Router) {
					r.Use(app.authMiddleware)
					r.Use(app.authorize(auth.ActionPostDelete, app.resourcePostFromCtx))
					r.Delete("/", app.deletePostHandler)
				})

				// Update post - requires auth + ownership
				r.Group(func(r chi.Router) {
					r.Use(app.authMiddleware)
					r.Use(app.authorize(auth.ActionPostUpdate, app.resourcePostFromCtx))
					r.Patch("/", app.updatePostHandler)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(app.userContextMiddleware)
			r.Get("/", app.getUserHandler)

			// Follow user - requires auth + policy check
			r.Group(func(r chi.Router) {
				r.Use(app.authMiddleware)
				r.Use(app.authorize(auth.ActionUserFollow, app.resourceUserFromCtx))
				r.Put("/{userID}/follow", app.followUserHandler)
			})

			// Unfollow user - requires auth + policy check
			r.Group(func(r chi.Router) {
				r.Use(app.authMiddleware)
				r.Use(app.authorize(auth.ActionUserUnfollow, app.resourceUserFromCtx))
				r.Delete("/{userID}/unfollow", app.unFollowUserHandler)
			})
		})

		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware)
			r.Get("/feed", app.getUserFeedHandler)
		})
	})

	return r
}

func (app *application) run(mux *chi.Mux) error {
	// Configure the HTTP server
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
	}

	// Create a channel to receive the shutdown operation result
	// This channel notifies the main goroutine whether graceful shutdown was successful or failed
	shutDown := make(chan error)

	// Start a separate goroutine for signal handling
	// This goroutine runs in the background and listens for system signals
	go func() {
		// Create a channel to catch SIGINT (Ctrl+C) and SIGTERM signals
		quit := make(chan os.Signal, 1)

		// Route these signals to the quit channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Wait for a signal to arrive (blocking operation)
		s := <-quit

		// Create a context with 5-second timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		// CRITICAL: Send the result of srv.Shutdown(ctx) to the shutDown channel
		// srv.Shutdown() attempts to gracefully close active connections
		// Returns nil if successful, error if failed
		// We send this result to the channel to notify the main goroutine
		shutDown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("Starting server", "address", app.config.addr, "environment", app.config.env)

	// Start the server (blocking operation)
	// Under normal conditions, this line runs continuously and never returns
	// Only returns if an error occurs or the server shuts down
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		// http.ErrServerClosed occurs during normal shutdown, this is not an error
		app.logger.Errorw("server error", "error", err)
	}

	// CRITICAL: Wait for the shutdown operation to complete
	// When a signal arrives, the goroutine above calls srv.Shutdown()
	// This line waits for that operation to finish (blocking operation)
	// The program stops here until shutdown is complete
	// This prevents the main goroutine from terminating prematurely
	err = <-shutDown

	if err != nil {
		app.logger.Errorw("server shutdown error", "error", err)
		return err
	}

	app.logger.Infow("Server gracefully stopped")
	return nil
}
