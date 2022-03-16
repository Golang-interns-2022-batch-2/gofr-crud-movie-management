package main

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	//"developer.zopsmart.com/go/gofr/examples/sample-https/handler"

	stores "github.com/anushi/newbatch/gofr-test/datastore/movie"
	delivery "github.com/anushi/newbatch/gofr-test/delivery/movie"
	services "github.com/anushi/newbatch/gofr-test/services/movie"
)

func main() {
	app := gofr.New()
	app.Server.ValidateHeaders = false

	//fmt.Println(app)
	s := stores.New()
	svc := services.New(s)
	h := delivery.New(svc)

	app.GET("/movie/{id}", h.GetByID)
	app.DELETE("/movie/{id}", h.Delete)
	app.PUT("/movie/{id}", h.Update)
	app.POST("/movie", h.Create)
	app.GET("/movie", h.GetAll)
	//	app.Server.HTTP.RedirectToHTTPS = true

	app.Start()
}
