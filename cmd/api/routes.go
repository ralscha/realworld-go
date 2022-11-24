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
		r.Group(func(r chi.Router) {
			r.Use(app.rwTransaction)
			r.Post("/users", app.usersRegistration)
		})
		r.Group(func(r chi.Router) {
			r.Use(app.readonlyTransaction)
			r.Post("/users/login", app.usersLogin)
			r.Get("/tags", app.tagsGet)
		})

		r.Group(func(r chi.Router) {
			r.Use(app.sessionManager.LoadAndSaveHeader)
			r.Use(app.readonlyTransaction)
			r.Get("/profiles/{username}", app.profilesGet)
			r.Get("/articles", app.articlesList)
			r.Get("/articles/{slug}", app.articleGet)
			r.Get("/articles/{slug}/comments", app.articlesGetComments)
		})
		r.Group(func(r chi.Router) {
			r.Use(app.sessionManager.LoadAndSaveHeader)
			r.Use(app.authenticatedOnly)

			r.Group(func(r chi.Router) {
				r.Use(app.rwTransaction)
				r.Put("/user", app.usersUpdate)
				r.Post("/profiles/{username}/follow", app.profilesFollow)
				r.Delete("/profiles/{username}/follow", app.profilesUnfollow)
				r.Post("/articles", app.articlesCreate)
				r.Put("/articles/{slug}", app.articlesUpdate)
				r.Delete("/articles/{slug}", app.articlesDelete)
				r.Post("/articles/{slug}/comments", app.articlesAddComment)
				r.Delete("/articles/{slug}/comments/:id", app.articlesDeleteComment)
				r.Post("/articles/{slug}/favorite", app.articlesFavorite)
				r.Delete("/articles/{slug}/favorite", app.articlesUnfavorite)
			})
			r.Group(func(r chi.Router) {
				r.Use(app.readonlyTransaction)
				r.Get("/user", app.usersGetCurrent)
				r.Get("/articles/feed", app.articlesFeed)
			})
		})
	})

	return mux
}
