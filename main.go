package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	h "github.com/shivam/Crud_Gofr/internal/http"
	"github.com/shivam/Crud_Gofr/internal/service"
	"github.com/shivam/Crud_Gofr/internal/store"
)

func main() {
	app := gofr.New()
	app.Server.ValidateHeaders = false

	MovieStore := store.New()
	MovieService := service.New(MovieStore)
	MovieHandler := h.New(MovieService)

	app.GET("/movie/{id}", MovieHandler.GetByID)
	app.GET("/movie", MovieHandler.GetAll)
	app.DELETE("/movie/{id}", MovieHandler.DeleteByID)
	app.POST("/movie", MovieHandler.Create)
	app.PUT("/movie/{id}", MovieHandler.UpdateByID)
	app.Start()
}
