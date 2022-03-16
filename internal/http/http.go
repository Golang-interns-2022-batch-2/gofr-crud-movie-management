package http

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/shivam/Crud_Gofr/internal/models"
	"github.com/shivam/Crud_Gofr/internal/service"
)

type response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
type data struct {
	Movies interface{} `json:"movie"`
}

type MovieHandler struct {
	service service.Interface
}

func New(h service.Interface) *MovieHandler {
	return &MovieHandler{service: h}
}

func (h *MovieHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	resp, err := h.service.GetByIDService(ctx, id)

	if err != nil {
		return nil, errors.EntityNotFound{
			Entity: "Movie",
			ID:     i,
		}
	}

	return resp, nil
}

func (h *MovieHandler) GetAll(ctx *gofr.Context) (interface{}, error) {
	movie, err := h.service.GetAllService(ctx)
	if err != nil {
		return nil, errors.Error("Internal Server Error")
	}

	result := response{
		Code:   200,
		Status: "Success",
		Data:   data{movie},
	}

	return result, nil
}

func (h *MovieHandler) DeleteByID(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	if err := h.service.DeleteService(ctx, id); err != nil {
		return nil, errors.Error("Internal Servor Error")
	}

	return "Deleted successfully", nil
}

func (h *MovieHandler) Create(ctx *gofr.Context) (interface{}, error) {
	var movie models.Movie
	if err := ctx.Bind(&movie); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	if movie.ID <= 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	resp, err := h.service.InsertService(ctx, &movie)
	if err != nil {
		return nil, errors.Error("Internal server Error")
	}

	return resp, nil
}

func (h *MovieHandler) UpdateByID(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")
	if i == "" {
		return nil, errors.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var movie models.Movie
	if err = ctx.Bind(&movie); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	movie.ID = id

	resp, err := h.service.UpdatedService(ctx, &movie)
	if err != nil {
		return nil, errors.Error("Internal Server Error")
	}

	return resp, nil
}
