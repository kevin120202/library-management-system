package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kevin120202/library-management-system/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/health", app.HealthCheck)

	r.Post("/api/users", app.UserHandler.HandleRegisterUser)
	r.Post("/api/authentication", app.TokenHandler.HandleCreateToken)

	r.Get("/api/books/{id}", app.BookHandler.HandleGetBookByID)
	r.Post("/api/books", app.BookHandler.HandleCreateBook)

	return r
}
