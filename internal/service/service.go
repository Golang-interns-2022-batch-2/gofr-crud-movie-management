package service

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/iamkakashi/movie-gofr/internal/model"
	"github.com/iamkakashi/movie-gofr/internal/store"
)

type MovieService struct {
	store store.MovieStorer
}

func New(ms store.MovieStorer) *MovieService {
	return &MovieService{ms}
}

func (service *MovieService) GetByID(ctx *gofr.Context, id int) (*model.Movie, error) {
	if id <= 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	movie, err := service.store.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return movie, nil
}
func (service *MovieService) Get(ctx *gofr.Context) ([]*model.Movie, error) {
	movies, err := service.store.Get(ctx)

	if err != nil {
		return nil, err
	}

	return movies, nil
}
func (service *MovieService) Create(ctx *gofr.Context, m *model.Movie) (*model.Movie, error) {
	if !Validation(m) {
		return nil, errors.InvalidParam{}
	}

	movie, err := service.store.Create(ctx, m)
	if err != nil {
		return nil, err
	}

	return movie, nil
}
func (service *MovieService) Update(ctx *gofr.Context, m *model.Movie) (*model.Movie, error) {
	if !Validation(m) {
		return nil, errors.InvalidParam{}
	}

	movie, err := service.store.Update(ctx, m)
	if err != nil {
		return nil, err
	}

	return movie, nil
}
func (service *MovieService) Delete(ctx *gofr.Context, id int) error {
	if id <= 0 {
		return errors.InvalidParam{Param: []string{"id"}}
	}

	err := service.store.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
