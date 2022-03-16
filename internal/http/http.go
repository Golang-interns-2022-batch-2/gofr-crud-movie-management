package http

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"github.com/shivam/Crud_Gofr/internal/models"
	"github.com/shivam/Crud_Gofr/internal/service"
)

type response struct {
	Data interface{} `json:"movie"`
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
	resp, err := h.service.GetByIDService(ctx, i)

	if err != nil {
		return nil, errors.EntityNotFound{
			Entity: "Movie",
			ID:     i,
		}
	}

	result := response{
		Data: resp,
	}

	return types.Response{Data: result}, nil
}

func (h *MovieHandler) GetAll(ctx *gofr.Context) (interface{}, error) {
	movie, err := h.service.GetAllService(ctx)
	if err != nil {
		return nil, errors.Error("Internal Server Error")
	}

	result := response{
		Data: movie,
	}

	return types.Response{Data: result}, nil
}

func (h *MovieHandler) DeleteByID(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")

	if err := h.service.DeleteService(ctx, i); err != nil {
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

	resp, err := h.service.InsertService(ctx, &movie)
	if err != nil {
		return nil, errors.Error("Internal server Error")
	}

	result := response{
		Data: resp,
	}

	return types.Response{Data: result}, nil
}

func (h *MovieHandler) UpdateByID(ctx *gofr.Context) (interface{}, error) {
	i := ctx.PathParam("id")

	var movie models.Movie

	if err := ctx.Bind(&movie); err != nil {
		ctx.Logger.Errorf("error in binding: %v", err)
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resp, err := h.service.UpdatedService(ctx, &movie, i)
	if err != nil {
		return nil, errors.Error("Internal Server Error")
	}

	result := response{
		Data: resp,
	}

	return types.Response{Data: result}, nil
}
