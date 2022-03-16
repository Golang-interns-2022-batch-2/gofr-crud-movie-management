package movie

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"golangprog/gofr-test/Internal/models"
	services "golangprog/gofr-test/Internal/service"
	"strconv"
)

type Handler struct {
	service services.Movie
}

func New(service services.Movie) *Handler {
	return &Handler{service: service}
}

type data struct {
	Movie interface{} `json:"movie"`
}

func (h *Handler) GetByID(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	body, err := h.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := types.Response{
		Data: data{
			body,
		},
	}

	return res, nil
}

func (h *Handler) GetAll(ctx *gofr.Context) (interface{}, error) {

	body, err := h.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	res := types.Response{
		Data: data{
			body,
		},
	}

	return res, nil
}

func (h Handler) Delete(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	err = h.service.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	res := types.Response{
		Data: data{
			"deleted successfully",
		},
	}
	return res, nil
}

func (h Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var movieRequestBody models.Movie
	err := ctx.Bind(&movieRequestBody)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	b, err := h.service.Create(ctx, &movieRequestBody)
	if err != nil {
		return nil, err
	}

	res := types.Response{
		Data: data{
			b,
		},
	}

	return res, nil
}

func (h Handler) Update(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	var movieObj models.Movie
	err = ctx.Bind(&movieObj)
	if err != nil {
		return nil, err
	}
	movieObj.ID = id
	body, err := h.service.Update(ctx, &movieObj)
	if err != nil {
		return nil, err
	}
	res := types.Response{
		Data: data{
			body,
		},
	}

	return res, nil
}
