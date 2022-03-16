//go:generate mockgen -destination=storeMock.go -package=store github.com/shivam/Crud_Gofr/internal/store MovieStorer
package store

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/shivam/Crud_Gofr/internal/models"
)

type MovieStorer interface {
	GetByID(ctx *gofr.Context, ID int) (*models.Movie, error)
	GetAll(ctx *gofr.Context) ([]*models.Movie, error)
	Create(ctx *gofr.Context, a *models.Movie) (*models.Movie, error)
	DeleteByID(ctx *gofr.Context, ID int) error
	Update(ctx *gofr.Context, ID int, updatedDetails *models.Movie) (*models.Movie, error)
}
