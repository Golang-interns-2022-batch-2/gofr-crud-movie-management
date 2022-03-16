package service

import (
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

func (service *MovieService) GetByIDService(ctx *gofr.Context, id int) (*models.Movie, error) {
	if id <= 0 {
		return nil, gofrerr.Error("id cannot be nagative or zero")
	}

	movie, err := service.store.GetByID(ctx, id)

	if err != nil {
		return nil, gofrerr.Error("error in service layer while calling getby id")
	}

	return movie, nil
}

func (service *MovieService) DeleteService(ctx *gofr.Context, id int) error {
	if id < 0 {
		return gofrerr.Error("id cannot be nagative")
	}

	_, check := service.store.GetByID(ctx, id)
	if check != nil {
		return gofrerr.Error("id no present in the database")
	}

	err := service.store.DeleteByID(ctx, id)

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

func (service *MovieService) UpdatedService(ctx *gofr.Context, a *models.Movie) (*models.Movie, error) {
	if a.ID < 0 {
		return nil, gofrerr.Error("ID cannot be negative")
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
	if a.Name == "" {
		return nil, gofrerr.Error("name cannot be null")
	}

	if a.Genre == "" {
		return nil, gofrerr.Error("genre cannot be null")
	}

	if a.Rating < 0 {
		return nil, gofrerr.Error("rating cannot be negative")
	}

	if a.ID <= 0 {
		return nil, gofrerr.Error("id cannot be negative")
	}

	if a.Plot == "" {
		return nil, gofrerr.Error("plot cannot be empty")
	}

	res, err := service.store.Create(ctx, a)

	if err != nil {
		return nil, gofrerr.Error("error in service layer while caaling create function")
	}

	return res, nil
}
