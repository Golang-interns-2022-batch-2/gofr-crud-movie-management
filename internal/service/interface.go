//go:generate mockgen -destination=serviceMock.go -package=service github.com/shivam/Crud_Gofr/internal/service Interface
package service

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/shivam/Crud_Gofr/internal/models"
)

type Interface interface {
	GetByIDService(ctx *gofr.Context, I string) (*models.Movie, error)
	GetAllService(ctx *gofr.Context) ([]*models.Movie, error)
	UpdatedService(ctx *gofr.Context, a *models.Movie, I string) (*models.Movie, error)
	InsertService(ctx *gofr.Context, a *models.Movie) (*models.Movie, error)
	DeleteService(ctx *gofr.Context, I string) error
}
