package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/model"
)

type MovieManager interface {
	InsertMovieService(*gofr.Context, *model.MovieModel) (*model.MovieModel, error)
	GetByIDService(*gofr.Context, int) (*model.MovieModel, error)
	DeleteByIDService(ctx *gofr.Context, id int) error
	UpdatedByIDService(*gofr.Context, *model.MovieModel) (*model.MovieModel, error)
	GetAllService(*gofr.Context) (*[]model.MovieModel, error)
}
