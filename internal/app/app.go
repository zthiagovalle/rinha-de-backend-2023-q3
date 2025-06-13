package app

import (
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/api"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/store"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/migrations"
)

type Application struct {
	DB            *pgxpool.Pool
	Logger        *log.Logger
	PersonHandler *api.PersonHandler
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigratePool(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	personStore := store.NewPostgresPersonStore(pgDB)

	personHandler := api.NewPersonHandler(logger, personStore)

	app := &Application{
		DB:            pgDB,
		Logger:        logger,
		PersonHandler: personHandler,
	}

	return app, nil
}
