package main

import (
	"net/http"
	"time"

	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/app"
	"github.com/zthiagovalle/rinha-de-backend-2023-q3/internal/routes"
)

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("server runnming on port 8080")

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
