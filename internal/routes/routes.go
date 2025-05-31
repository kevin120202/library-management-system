package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kevin120202/library-management-system/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Get("/api/books/{id}", app.Middleware.RequireUser(app.BookHandler.HandleGetBookByID))
		r.Post("/api/books", app.Middleware.RequireUser(app.BookHandler.HandleCreateBook))

		r.Post("/api/logout", app.UserHandler.HandleLogoutUser)
	})

	r.Get("/api/health", app.HealthCheck)

	r.Post("/api/users", app.UserHandler.HandleRegisterUser)
	r.Post("/api/authentication", app.TokenHandler.HandleCreateToken)

	return r
}
