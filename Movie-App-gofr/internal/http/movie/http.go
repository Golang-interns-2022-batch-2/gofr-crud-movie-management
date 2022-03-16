package movie

import (
	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/types"

	"github.com/RicheshZopsmart/Movie-App-gofr/internal/model"
	serv "github.com/RicheshZopsmart/Movie-App-gofr/internal/service"
)

// For response Formatting
type DataField struct {
	Movie *model.MovieModel `json:"movie"`
}

type DataList struct {
	Movie *[]model.MovieModel `json:"movie"`
}

type GetAllModelResponse struct {
	Code   int      `json:"code"`
	Status string   `json:"status"`
	Data   DataList `json:"data"`
}

type ModelResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
type ResponseModel struct {
	Code    int
	Status  string
	Message string
}
type Handler struct {
	ServiceHandler serv.MovieManager
}

func New(movieinterface serv.MovieManager) *Handler {
	return &Handler{ServiceHandler: movieinterface}
}

func (h Handler) GetByIDRequest(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	mObj, err := h.ServiceHandler.GetByIDService(ctx, id)

	if err != nil {
		return nil, &errors.Response{StatusCode: 500, Code: "500", Reason: "Internal Server Error"}
	}

	movies := DataField{
		Movie: mObj,
	}
	respObj := ModelResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   movies,
	}

	return respObj, nil
}

func (h Handler) CreateMovieRequest(ctx *gofr.Context) (interface{}, error) {
	var movieObj model.MovieModel

	err := ctx.Bind(&movieObj)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resMovieObj, err := h.ServiceHandler.InsertMovieService(ctx, &movieObj)

	if err != nil {
		return nil, errors.Error("Internal Server Error")
	}

	movies := DataField{
		Movie: resMovieObj,
	}

	respObj := ModelResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   movies,
	}

	return respObj, nil
}

func (h Handler) DeleteByIDRequest(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err = h.ServiceHandler.DeleteByIDService(ctx, id)

	if err != nil {
		return nil, err
	}

	response := types.Response{
		Data: ResponseModel{Code: 200, Status: "SUCCESS", Message: "Movie Deleted Successfully"},
	}

	return response, nil
}

func (h Handler) UpdateByIDRequest(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var resMovieObj model.MovieModel

	err = ctx.Bind(&resMovieObj)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resMovieObj.ID = id
	newResMovieObj, err := h.ServiceHandler.UpdatedByIDService(ctx, &resMovieObj)

	if err != nil {
		return nil, errors.Error("internal server error")
	}

	response := types.Response{
		Data: DataField{Movie: newResMovieObj},
	}

	return response, err
}

func (h Handler) GetAllRequest(ctx *gofr.Context) (interface{}, error) {
	movies, err := h.ServiceHandler.GetAllService(ctx)

	if err != nil {
		return nil, errors.Error("internal server error")
	}

	moviesObj := DataList{
		Movie: movies,
	}

	response := types.Response{
		Data: moviesObj,
	}

	return response, nil
}
