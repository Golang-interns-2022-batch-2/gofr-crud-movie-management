package movie

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"
	"fmt"
	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/models"
	"log"
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	//"github.com/anushi/newbatch/gofr-test/models"
	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/services"
)

type Handler struct {
	service services.Movie
}

func New(c services.Movie) *Handler {
	return &Handler{service: c}
}

func ErrorResponse(statusCode int, msg string) []byte {
	errorMsg := fmt.Sprintf(`{"code":%v, "status":%v, "message":%v"}`, statusCode, "Error", msg)
	return []byte(errorMsg)
}

type data struct {
	Movie interface{} `json:"movie"`
}
type response struct {
	Data interface{} `json:"movie"`
}

func (handler *Handler) GetByID(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	resp, err := handler.service.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	response := response{
		Data: data{
			resp},
	}

	return types.Response{Data: response}, nil

}

func (handler *Handler) Delete(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		log.Println("Error in converting id")
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err = handler.service.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	response := types.Response{
		Data: "Movie Deleted Successfully",
	}
	return response, nil
}

func (handler *Handler) Update(ctx *gofr.Context) (interface{}, error) {

	i := ctx.PathParam("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	var movie models.Movie
	err = ctx.Bind(&movie)
	if err != nil {
		return nil, err
	}
	movie.ID = id
	body, err := handler.service.Update(ctx, &movie)
	if err != nil {
		return nil, err
	}

	response := response{
		Data: data{
			body,
		},
	}

	return types.Response{Data: response}, nil
}

func (handler *Handler) Create(ctx *gofr.Context) (interface{}, error) {

	var movie *models.Movie
	err := ctx.Bind(&movie)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	mov, err := handler.service.Create(ctx, movie)
	if err != nil {
		return nil, err
	}
	response := response{
		Data: data{mov},
	}
	return types.Response{Data: response}, nil
}

func (handler *Handler) GetAll(ctx *gofr.Context) (interface{}, error) {

	//var movie []*models.Movie
	resp, err := handler.service.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	//body, _ := json.Marshal(resp)
	//
	//msg := fmt.Sprintf(`{"code":%v, "status":%v, "data":{"movie":{%v"}}}`, http.StatusOK, "Success", string(body))
	//
	//return msg, nil

	response := response{
		Data: data{resp},
	}

	return types.Response{Data: response}, nil

}
