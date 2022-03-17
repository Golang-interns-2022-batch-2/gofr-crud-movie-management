package movie

import (
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/errors"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/datastore"
	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/models"
)

type Service struct {
	movStore datastore.Movie
}

func New(m datastore.Movie) *Service {
	return &Service{movStore: m}
}

func (s *Service) GetByID(ctx *gofr.Context, id int) (*models.Movie, error) {
	if id < 1 {
		return nil, errors.InvalidParam{Param: []string{"id"}}

	}

	return s.movStore.GetByID(ctx, id)
}
func (s *Service) Delete(ctx *gofr.Context, id int) error {
	if id <= 0 {
		return errors.InvalidParam{Param: []string{"id"}}
	}
	err := s.movStore.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(ctx *gofr.Context, mov *models.Movie) (*models.Movie, error) {
	if mov.ID < 1 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}
	return s.movStore.Update(ctx, mov)
}

func (s *Service) Create(ctx *gofr.Context, mov *models.Movie) (*models.Movie, error) {
	if mov.ID < 0 {
		return nil, errors.InvalidParam{}

	}
	if mov.Name == "" {
		return nil, errors.InvalidParam{}

	}
	if mov.Genre == "" {
		return nil, errors.InvalidParam{}

	}

	res, err := s.movStore.Create(ctx, mov)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}

	}

	return res, nil
}

func (s *Service) GetAll(ctx *gofr.Context) ([]*models.Movie, error) {
	i, err := s.movStore.GetAll(ctx)
	if err != nil {
		return nil, sql.ErrNoRows
	}

	return i, nil
}
