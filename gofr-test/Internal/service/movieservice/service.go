package movieservice

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"errors"
	"golangprog/gofr-test/Internal/models"
	stores "golangprog/gofr-test/Internal/store"
)

type Service struct {
	store stores.Movie
}

func New(store stores.Movie) *Service {
	return &Service{store}
}

func (se *Service) GetByID(ctx *gofr.Context, id int) (*models.Movie, error) {
	if id <= 0 {
		return nil, errors.New("error invalid id")
	}

	movieObj, err := se.store.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return movieObj, nil
}

func (se *Service) GetAll(ctx *gofr.Context) ([]*models.Movie, error) {
	movieObj, err := se.store.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return movieObj, nil
}

func (se *Service) Delete(ctx *gofr.Context, id int) error {
	err := se.store.Delete(ctx, id)

	if err != nil {
		return errors.New("error invalid id")
	}

	return nil
}

func (se *Service) Create(ctx *gofr.Context, movieObj *models.Movie) (*models.Movie, error) {
	if movieObj.Name == "" {
		return nil, errors.New("error invalid name")
	}

	if movieObj.Plot == "" {
		return nil, errors.New("error invalid plot")
	}

	movieObjs, err := se.store.Create(ctx, movieObj)

	if err != nil {
		return nil, err
	}

	return movieObjs, nil
}

func (se *Service) Update(ctx *gofr.Context, movieObj *models.Movie) (*models.Movie, error) {
	if movieObj.ID < 0 {
		return nil, errors.New("error invalid id")
	}

	if movieObj.Name == "" {
		return nil, errors.New("error invalid name")
	}

	if movieObj.Plot == "" {
		return nil, errors.New("error invalid plot")
	}

	movieObjs, err := se.store.Update(ctx, movieObj)

	if err != nil {
		return nil, err
	}

	return movieObjs, nil
}
