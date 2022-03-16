package movie

import (
	"database/sql"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/model"
)

const format = "2006-01-02 15:04:05"

type DBStore struct {
}

func NewDBHandler() *DBStore {
	return &DBStore{}
}

func (s DBStore) CreateMovie(ctx *gofr.Context, movieObj *model.MovieModel) (*model.MovieModel, error) {
	query := "insert into movie_details(name,genre,rating,plot,released,releaseDate) values(?,?,?,?,?,?); "

	res2, execErr := ctx.DB().ExecContext(
		ctx, query, movieObj.Name,
		movieObj.Genre, movieObj.Rating, movieObj.Plot,
		movieObj.Released, movieObj.ReleaseDate,
	)

	if execErr != nil {
		return nil, errors.Error("Internal Server Error")
	}

	resultedID, _ := res2.LastInsertId()
	resMObj, err := s.GetByID(ctx, int(resultedID))

	return resMObj, err
}
func (s DBStore) GetByID(ctx *gofr.Context, id int) (*model.MovieModel, error) {
	query := "select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from movie_details where id = ? and deletedAt is null;"

	var movie model.MovieModel

	var releaseDateScan sql.NullString

	var CreatedAtScan sql.NullString

	err := ctx.DB().QueryRowContext(ctx, query, id).Scan(
		&movie.ID, &movie.Name, &movie.Genre,
		&movie.Rating, &releaseDateScan, &movie.UpdatedAt,
		&CreatedAtScan, &movie.Plot, &movie.Released,
	)

	if err != nil {
		return nil, err
	}

	movie.ReleaseDate = releaseDateScan.String
	movie.CreatedAt = CreatedAtScan.String

	return &movie, nil
}
func (s DBStore) UpdateByID(ctx *gofr.Context, movieObj *model.MovieModel) (*model.MovieModel, error) {
	query := "update movie_details set rating=?,plot=?,releaseDate=?,updatedAt=? where deletedAt is null and id=?;"

	updateTime := time.Now()

	movieObj.UpdatedAt = time.Now().Format(format)

	_, err := ctx.DB().ExecContext(ctx, query, movieObj.Rating, movieObj.Plot, movieObj.ReleaseDate, updateTime, movieObj.ID)

	return movieObj, err
}
func (s DBStore) DeleteByID(ctx *gofr.Context, id int) error {
	query := "update movie_details set deletedAt = ? where id = ? and deletedAt is null;"

	updateTime := time.Now().Format(format)

	_, err := ctx.DB().ExecContext(ctx, query, updateTime, id)

	return err
}
func (s DBStore) GetAll(ctx *gofr.Context) (*[]model.MovieModel, error) {
	query := "select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from movie_details where deletedAt is null;"

	var mObj []model.MovieModel

	var tmpObj model.MovieModel

	rows, err := ctx.DB().QueryContext(ctx, query)

	if err != nil {
		return nil, errors.Error("Couldn't execute query")
	}

	for rows.Next() {
		err := rows.Scan(
			&tmpObj.ID,
			&tmpObj.Name,
			&tmpObj.Genre,
			&tmpObj.Rating,
			&tmpObj.ReleaseDate,
			&tmpObj.UpdatedAt,
			&tmpObj.CreatedAt,
			&tmpObj.Plot,
			&tmpObj.Released,
		)

		if err != nil {
			return nil, errors.Error("Scan Error")
		}

		mObj = append(mObj, tmpObj)
	}

	return &mObj, nil
}
