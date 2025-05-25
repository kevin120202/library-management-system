package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kevin120202/library-management-system/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/api/health", app.HealthCheck)

	r.Get("/api/book/{id}", app.BookHandler.HandleGetBookByID)
	r.Post("/api/book", app.BookHandler.HandleCreateBook)

	return r
}
