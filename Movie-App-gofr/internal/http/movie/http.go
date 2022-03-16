package movie

import (
	"database/sql"

	"strconv"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
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
	Code   int       `json:"code"`
	Status string    `json:"status"`
	Data   DataField `json:"data"`
}
type ResponseModel struct {
	Code    int
	Status  string
	Message string
}
type Handler struct {
	ServiceHandler serv.MovieServiceInterface
}

func New(movieinterface serv.MovieServiceInterface) *Handler {
	return &Handler{ServiceHandler: movieinterface}
}

func jsonSuccessMessage() ResponseModel {
	ErrObj := ResponseModel{Code: 200, Message: "Successfully Executed!", Status: "Success"}
	return ErrObj
}

func (mh Handler) GetByIDRequest(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	mObj, err := mh.ServiceHandler.GetByIDService(ctx, id)

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "movies", ID: ctx.PathParam("id")}
	}

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

func (mh Handler) CreateMovieRequest(ctx *gofr.Context) (interface{}, error) {
	var movieObj model.MovieModel

	err := ctx.Bind(&movieObj)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	resMovieObj, err := mh.ServiceHandler.InsertMovieService(ctx, &movieObj)

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

func (mh Handler) DeleteByIDRequest(ctx *gofr.Context) (interface{}, error) {
	id, err := strconv.Atoi(ctx.PathParam("id"))
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err = mh.ServiceHandler.DeleteByIDService(ctx, id)
	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "movie", ID: strconv.Itoa(id)}
	}

	if err != nil {
		return nil, err
	}

	return jsonSuccessMessage(), nil
}

func (mh Handler) UpdateByIDRequest(ctx *gofr.Context) (interface{}, error) {
	rawID := ctx.PathParam("id")
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
	newResMovieObj, err := mh.ServiceHandler.UpdatedByIDService(ctx, &resMovieObj)

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "movie", ID: rawID}
	}

	if err != nil {
		return nil, errors.Error("internal server error")
	}

	movies := DataField{
		Movie: newResMovieObj,
	}

	respObj := ModelResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   movies,
	}

	return respObj, err
}

func (mh Handler) GetAllRequest(ctx *gofr.Context) (interface{}, error) {
	movies, err := mh.ServiceHandler.GetAllService(ctx)

	if err != nil {
		return nil, errors.Error("internal server error")
	}

	moviesObj := DataList{
		Movie: movies,
	}

	respObj := GetAllModelResponse{
		Code:   200,
		Status: "SUCCESS",
		Data:   moviesObj,
	}

	return respObj, nil
}
