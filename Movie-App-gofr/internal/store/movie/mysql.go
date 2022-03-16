package movie

import (
	"database/sql"
	"strconv"
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

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "movies", ID: ctx.PathParam("id")}
	}

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

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "movie", ID: ctx.PathParam("id")}
	}

	return movieObj, err
}
func (s DBStore) DeleteByID(ctx *gofr.Context, id int) error {
	query := "update movie_details set deletedAt = ? where id = ? and deletedAt is null;"

	updateTime := time.Now().Format(format)

	_, err := ctx.DB().ExecContext(ctx, query, updateTime, id)

	if err == sql.ErrNoRows {
		return errors.EntityNotFound{Entity: "movie", ID: strconv.Itoa(id)}
	}

	return err
}
func (s DBStore) GetAll(ctx *gofr.Context) (*[]model.MovieModel, error) {
	query := "select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from movie_details where deletedAt is null;"

	var movies []model.MovieModel

	var movie model.MovieModel

	rows, err := ctx.DB().QueryContext(ctx, query)

	if err != nil {
		return nil, errors.Error("Couldn't execute query")
	}

	for rows.Next() {
		err := rows.Scan(
			&movie.ID,
			&movie.Name,
			&movie.Genre,
			&movie.Rating,
			&movie.ReleaseDate,
			&movie.UpdatedAt,
			&movie.CreatedAt,
			&movie.Plot,
			&movie.Released,
		)

		if err != nil {
			return nil, errors.Error("Scan Error")
		}

		movies = append(movies, movie)
	}

	return &movies, nil
}
