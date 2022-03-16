package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	handler "github.com/iamkakashi/movie-gofr/internal/http"
	"github.com/iamkakashi/movie-gofr/internal/service"
	"github.com/iamkakashi/movie-gofr/internal/store"
)

func main() {
	app := gofr.New()
	app.Server.ValidateHeaders = false
	movieStore := store.New()
	movieService := service.New(movieStore)
	movieHandler := handler.New(movieService)
	app.GET("/movies", movieHandler.Get)
	app.GET("/movies/{id}", movieHandler.GetByID)
	app.POST("/movies", movieHandler.Post)
	app.PUT("/movies/{id}", movieHandler.Put)
	app.DELETE("/movies/{id}", movieHandler.Delete)
	app.Start()
}
