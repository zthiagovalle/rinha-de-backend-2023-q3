package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/api"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/store"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/migrations"
)

type Application struct {
	DB            *sql.DB
	Logger        *log.Logger
	PersonHandler *api.PersonHandler
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

	personHandler := api.NewPersonHandler(logger)

	app := &Application{
		DB:            pgDB,
		Logger:        logger,
		PersonHandler: personHandler,
	}

	return app, nil
}
