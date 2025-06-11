package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/contagem-pessoas", app.PersonHandler.HandleCountPersons)
	r.Get("/pessoas/{id}", app.PersonHandler.HandleGetPerson)
	r.Get("/pessoas/", app.PersonHandler.HandleGetPerson)
	r.Post("/pessoas", app.PersonHandler.HandleCreatePerson)

	return r
}
