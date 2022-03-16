package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/model"
)

type MovieInterface interface {
	DeleteByID(ctx *gofr.Context, id int) error
	CreateMovie(*gofr.Context, *model.MovieModel) (*model.MovieModel, error)
	GetByID(ctx *gofr.Context, id int) (*model.MovieModel, error)
	UpdateByID(ctx *gofr.Context, movieObj *model.MovieModel) (*model.MovieModel, error)
	GetAll(*gofr.Context) (*[]model.MovieModel, error)
}
