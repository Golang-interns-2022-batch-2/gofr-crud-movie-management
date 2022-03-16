package main

import (

	// Own Packages
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	ht "github.com/RicheshZopsmart/Movie-App-gofr/internal/http/movie"
	serv "github.com/RicheshZopsmart/Movie-App-gofr/internal/service/movie"

	"github.com/RicheshZopsmart/Movie-App-gofr/internal/store/movie"
)

func main() {
	datastore := movie.NewDBHandler()
	servicestore := serv.NewMovieServiceHandler(datastore)
	handler := ht.New(servicestore)

	app := gofr.New()
	app.Server.ValidateHeaders = false

	app.GET("/movies/{id}", handler.GetByIDRequest)
	app.GET("/movies", handler.GetAllRequest)

	app.POST("/movies", handler.CreateMovieRequest)
	app.PUT("/movies/{id}", handler.UpdateByIDRequest)

	app.DELETE("/movies/{id}", handler.DeleteByIDRequest)
	app.Start()
}
