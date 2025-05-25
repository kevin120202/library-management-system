package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kevin120202/library-management-system/internal/api"
)

type Application struct {
	Logger      *log.Logger
	BookHandler *api.BookHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	bookHandler := api.NewBookHandler(logger)

	app := &Application{
		Logger:      logger,
		BookHandler: bookHandler,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
