package datastore

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/anushi/newbatch/gofr-test/models"
)

type Movie interface {
	GetByID(ctx *gofr.Context, id int) (*models.Movie, error)
	Delete(ctx *gofr.Context, id int) error
	Update(ctx *gofr.Context, mov *models.Movie) (*models.Movie, error)
	Create(ctx *gofr.Context, mov *models.Movie) (*models.Movie, error)
	GetAll(ctx *gofr.Context) ([]*models.Movie, error)
}
