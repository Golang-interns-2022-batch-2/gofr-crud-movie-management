//go:generate mockgen -destination=mock_interface.go -package=service github.com/iamkakashi/movie-gofr/internal/service MovieServicer

package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/iamkakashi/movie-gofr/internal/model"
)

type MovieServicer interface {
	GetByID(ctx *gofr.Context, id int) (*model.Movie, error)
	Get(ctx *gofr.Context) ([]*model.Movie, error)
	Create(ctx *gofr.Context, m *model.Movie) (*model.Movie, error)
	Update(ctx *gofr.Context, m *model.Movie) (*model.Movie, error)
	Delete(ctx *gofr.Context, id int) error
}
