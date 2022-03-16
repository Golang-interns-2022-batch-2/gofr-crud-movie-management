package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"golangprog/gofr-test/Internal/models"
)

type Movie interface {
	GetByID(ctx *gofr.Context, id int) (*models.Movie, error)
	GetAll(ctx *gofr.Context) ([]*models.Movie, error)
	Delete(ctx *gofr.Context, id int) error
	Create(ctx *gofr.Context, movieObj *models.Movie) (*models.Movie, error)
	Update(ctx *gofr.Context, movieObj *models.Movie) (*models.Movie, error)
}
