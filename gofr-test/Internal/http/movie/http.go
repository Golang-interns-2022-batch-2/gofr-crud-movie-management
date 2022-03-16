package movie

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"errors"
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

type response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
type data struct {
	Movie interface{} `json:"movie"`
}

func (h *Handler) GetByID(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, err
	}

	body, err := h.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	response := response{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{body},
	}

	return response, nil
}

func (h *Handler) GetAll(ctx *gofr.Context) (interface{}, error) {

	body, err := h.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	response := response{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{body},
	}

	return response, nil
}

func (h Handler) Delete(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, err
	}
	err = h.service.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return "Deleted successfully", nil
}

func (h Handler) Create(ctx *gofr.Context) (interface{}, error) {
	var movieRequestBody models.Movie
	err := ctx.Bind(&movieRequestBody)
	if err != nil {
		return nil, err
	}

	b, err := h.service.Create(ctx, &movieRequestBody)
	if err != nil {
		return nil, errors.New("Error in inserting")
	}

	response := response{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{b},
	}

	return response, nil
}

func (h Handler) Update(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, err
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
	response := response{
		Code:   200,
		Status: "SUCCESS",
		Data:   data{body},
	}

	return response, nil
}
