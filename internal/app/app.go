package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kevin120202/library-management-system/internal/api"
	"github.com/kevin120202/library-management-system/internal/middleware"
	"github.com/kevin120202/library-management-system/internal/store"
	"github.com/kevin120202/library-management-system/migrations"
)

type Application struct {
	Logger       *log.Logger
	UserHandler  *api.UserHandler
	TokenHandler *api.TokenHandler
	Middleware   middleware.UserMiddleware
	BookHandler  *api.BookHandler
	DB           *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	userHandler := api.NewUserHandler(userStore, tokenStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	bookHandler := api.NewBookHandler(logger)

	middlewareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:       logger,
		UserHandler:  userHandler,
		TokenHandler: tokenHandler,
		Middleware:   middlewareHandler,
		BookHandler:  bookHandler,
		DB:           pgDB,
	}

	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
