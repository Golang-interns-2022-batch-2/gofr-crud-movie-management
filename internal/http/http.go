package http

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/iamkakashi/movie-gofr/internal/model"
	"github.com/iamkakashi/movie-gofr/internal/service"
)

type MovieHandler struct {
	service service.MovieServicer
}

func New(ms service.MovieServicer) *MovieHandler {
	return &MovieHandler{ms}
}

type apiResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
type data struct {
	Movies interface{} `json:"movie"`
}

func (handler *MovieHandler) Get(ctx *gofr.Context) (interface{}, error) {
	movie, err := handler.service.Get(ctx)
	if err != nil {
		return nil, err
	}

	response := apiResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{movie},
	}

	return response, nil
}

func (handler *MovieHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	movie, err := handler.service.GetByID(ctx, id)
	if err != nil {
		return nil, errors.EntityNotFound{
			Entity: "movie",
			ID:     i,
		}
	}

	response := apiResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{movie},
	}

	return response, nil
}

func (handler *MovieHandler) Post(ctx *gofr.Context) (interface{}, error) {
	var movie model.Movie
	if err := ctx.Bind(&movie); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := handler.service.Create(ctx, &movie)

	if err != nil {
		return nil, err
	}

	response := apiResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{resp},
	}

	return response, nil
}
func (handler *MovieHandler) Put(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var movie model.Movie
	if err = ctx.Bind(&movie); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	movie.ID = id

	resp, err := handler.service.Update(ctx, &movie)
	if err != nil {
		return nil, err
	}

	response := apiResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{resp},
	}

	return response, nil
}

func (handler *MovieHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	if err := handler.service.Delete(ctx, id); err != nil {
		return nil, err
	}

	return "Deleted successfully", nil
}
