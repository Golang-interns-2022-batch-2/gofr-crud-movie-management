package main

import (
	//"developer.zopsmart.com/go/gofr/examples/sample-https/handler"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"golangprog/gofr-test/Internal/http/movie"
	"golangprog/gofr-test/Internal/service/movieservice"
	"golangprog/gofr-test/Internal/store/moviestore"
)

func main() {
	app := gofr.New()
	app.Server.ValidateHeaders = false
	s := moviestore.New()
	se := movieservice.New(s)
	h := movie.New(se)
	app.GET("/movie/{id}", h.GetByID)
	app.GET("/movie", h.GetAll)
	app.DELETE("/movie/{id}", h.Delete)
	app.POST("/movie", h.Create)
	app.PUT("/movie/{id}", h.Update)

	//app.Server.HTTP.RedirectToHTTPS = true

	app.Start()
}
