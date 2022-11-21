package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"net/http"
	"realworldgo.rasc.ch/internal/config"
	"realworldgo.rasc.ch/internal/response"
	"time"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(response.NotFound)
	mux.MethodNotAllowed(response.MethodNotAllowed)

	// Middleware
	mux.Use(middleware.RealIP)
	if app.config.Environment == config.Development {
		mux.Use(middleware.Logger)
	}

	mux.Use(middleware.Recoverer)
	mux.Use(httprate.LimitAll(1_000, 1*time.Minute))
	mux.Use(middleware.Timeout(15 * time.Second))
	mux.Use(middleware.NoCache)

	mux.Route("/api", func(r chi.Router) {
		r.Post("/users/login", app.usersLogin)
		r.Post("/users", app.usersRegistration)
		r.Mount("/", app.authenticatedRouter())
	})

	return mux
}

func (app *application) authenticatedRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(app.sessionManager.LoadAndSaveHeader)
	r.Use(app.authenticatedOnly)

	return r
}

func (app *application) authenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := app.sessionManager.Get(r.Context(), "userID")
		userID, ok := value.(int64)
		if ok && userID > 0 {
			next.ServeHTTP(w, r)
		} else {
			response.Forbidden(w)
		}
	})
}
