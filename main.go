package main

import (
	"fmt"
	"net/http"
	"os"
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
		Addr:         fmt.Sprintf(":%s", os.Getenv("API_PORT")),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("server running on port %s", os.Getenv("API_PORT"))

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
