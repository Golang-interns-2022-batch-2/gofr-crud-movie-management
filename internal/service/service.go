package service

import (
	"strconv"

	gofrerr "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/shivam/Crud_Gofr/internal/models"
	"github.com/shivam/Crud_Gofr/internal/store"
)

type MovieService struct {
	store store.MovieStorer
}

func New(s store.MovieStorer) *MovieService {
	return &MovieService{s}
}

func (service *MovieService) GetByIDService(ctx *gofr.Context, id string) (*models.Movie, error) {
	if id == "" {
		return nil, gofrerr.MissingParam{Param: []string{"id"}}
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, gofrerr.InvalidParam{Param: []string{"id"}}
	}

	if i <= 0 {
		return nil, gofrerr.InvalidParam{Param: []string{"id"}}
	}

	movie, err := service.store.GetByID(ctx, i)

	if err != nil {
		return nil, gofrerr.Error("error in service layer while calling getby id")
	}

	return movie, nil
}

func (service *MovieService) DeleteService(ctx *gofr.Context, i string) error {
	if i == "" {
		return gofrerr.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return gofrerr.InvalidParam{Param: []string{"id"}}
	}

	if id < 0 {
		return gofrerr.Error("id cannot be nagative")
	}

	_, check := service.store.GetByID(ctx, id)
	if check != nil {
		return gofrerr.Error("id no present in the database")
	}

	err = service.store.DeleteByID(ctx, id)

	if err != nil {
		return gofrerr.Error("error in service layer for deleting id")
	}

	return nil
}

func (service *MovieService) GetAllService(ctx *gofr.Context) ([]*models.Movie, error) {
	res, err := service.store.GetAll(ctx)
	if err != nil {
		return nil, gofrerr.Error("error1 in service layer while calling GetAll function")
	}

	return res, nil
}

func (service *MovieService) UpdatedService(ctx *gofr.Context, a *models.Movie, i string) (*models.Movie, error) {
	if i == "" {
		return nil, gofrerr.MissingParam{Param: []string{"id"}}
	}

	id, err := strconv.Atoi(i)
	if err != nil {
		return nil, gofrerr.InvalidParam{Param: []string{"id"}}
	}

	a.ID = id

	if a.ID < 0 {
		return nil, gofrerr.InvalidParam{Param: []string{"id"}}
	}

	_, check := service.store.GetByID(ctx, a.ID)
	if check != nil {
		return nil, gofrerr.Error("id no present in the database")
	}

	movie, err := service.store.Update(ctx, a.ID, a)

	if err != nil {
		return nil, gofrerr.Error("error while calling update in service layer")
	}

	return movie, nil
}

func (service *MovieService) InsertService(ctx *gofr.Context, a *models.Movie) (*models.Movie, error) {
	if a.ID <= 0 {
		return nil, gofrerr.InvalidParam{Param: []string{"id"}}
	}

	if a.Name == "" {
		return nil, gofrerr.InvalidParam{Param: []string{"name"}}
	}

	if a.Genre == "" {
		return nil, gofrerr.InvalidParam{Param: []string{"genre"}}
	}

	if a.Rating < 0 {
		return nil, gofrerr.InvalidParam{Param: []string{"rating"}}
	}

	if a.Plot == "" {
		return nil, gofrerr.InvalidParam{Param: []string{"plot"}}
	}

	res, err := service.store.Create(ctx, a)

	if err != nil {
		return nil, gofrerr.Error("error in service layer while caaling create function")
	}

	return res, nil
}
