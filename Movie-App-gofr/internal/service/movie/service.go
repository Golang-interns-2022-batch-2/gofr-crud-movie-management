package movie

import (
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/model"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/store"
)

type ServiceHandler struct {
	movServHandler store.MovieInterface
}

func NewMovieServiceHandler(storeInterface store.MovieInterface) *ServiceHandler {
	return &ServiceHandler{movServHandler: storeInterface}
}

func (moviehandler ServiceHandler) InsertMovieService(ctx *gofr.Context, movieObj *model.MovieModel) (*model.MovieModel, error) {
	name := movieObj.Name
	genre := movieObj.Genre

	nameErr := isNameEmpty(name)

	if nameErr {
		return nil, &errors.Response{
			StatusCode: 400,
			Code:       "INVLD_NAME",
			Reason:     "name is empty",
		}
	}

	genreErr := isGenreEmpty(genre)

	if genreErr {
		return nil, &errors.Response{
			StatusCode: 400,
			Code:       "INVLD_GENRE",
			Reason:     "genre is empty",
		}
	}

	res, err := moviehandler.movServHandler.CreateMovie(ctx, movieObj)

	return res, err
}

func (moviehandler ServiceHandler) GetByIDService(ctx *gofr.Context, id int) (*model.MovieModel, error) {
	if id < 0 {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	mObj, err := moviehandler.movServHandler.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return mObj, nil
}
func (moviehandler ServiceHandler) GetAllService(ctx *gofr.Context) (*[]model.MovieModel, error) {

	resultObj, err := moviehandler.movServHandler.GetAll(ctx)

	return resultObj, err
}

func (moviehandler ServiceHandler) DeleteByIDService(ctx *gofr.Context, id int) error {
	if id < 0 {
		return errors.Error("negative ID found")
	}

	_, err := moviehandler.movServHandler.GetByID(ctx, id)

	if err != nil {
		return err
	}

	err = moviehandler.movServHandler.DeleteByID(ctx, id)

	return err
}

func (moviehandler ServiceHandler) UpdatedByIDService(ctx *gofr.Context, mObj *model.MovieModel) (*model.MovieModel, error) {
	id := mObj.ID

	if id < 0 {
		return nil, errors.Error("negative ID found")
	}

	getMovObj, err := moviehandler.movServHandler.GetByID(ctx, mObj.ID)

	if getMovObj == nil {
		return nil, err
	}

	resultObj, err := moviehandler.movServHandler.UpdateByID(ctx, mObj)

	if err != nil {
		return nil, err
	}

	resultObj.Name = getMovObj.Name
	resultObj.Genre = getMovObj.Genre
	resultObj.CreatedAt = getMovObj.CreatedAt

	return resultObj, nil
}

// func (moviehandler ServiceHandler) GetAllService() (*[]model.MovieModel, error) {
// 	resultObj, err := moviehandler.movServHandler.GetAll()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return resultObj, nil
// }
