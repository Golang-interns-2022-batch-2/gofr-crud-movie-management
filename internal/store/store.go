package store

import (
	"fmt"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/iamkakashi/movie-gofr/internal/model"
)

type MovieStore struct {
}

func New() *MovieStore {
	return &MovieStore{}
}

func (dbstore *MovieStore) GetByID(ctx *gofr.Context, id int) (*model.Movie, error) {
	row := ctx.DB().QueryRowContext(ctx,
		"SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL AND ID = ?", id,
	)
	movie := model.Movie{}
	err := row.Scan(
		&movie.ID,
		&movie.Name,
		&movie.Genre,
		&movie.Rating,
		&movie.ReleaseDate,
		&movie.UpdatedAt,
		&movie.CreatedAt,
		&movie.Plot,
		&movie.Released)

	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, errors.EntityNotFound{Entity: "movie", ID: fmt.Sprint(id)}
	}

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return &movie, nil
}

func (dbstore *MovieStore) Get(ctx *gofr.Context) ([]*model.Movie, error) {
	rows, err := ctx.DB().QueryContext(ctx,
		"SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL",
	)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	movies := []*model.Movie{}
	numRows := 0

	defer rows.Close()

	for rows.Next() {
		movie := model.Movie{}
		err = rows.Scan(
			&movie.ID,
			&movie.Name,
			&movie.Genre,
			&movie.Rating,
			&movie.ReleaseDate,
			&movie.UpdatedAt,
			&movie.CreatedAt,
			&movie.Plot,
			&movie.Released)

		if err != nil {
			return nil, errors.DB{Err: err}
		}

		movies = append(movies, &movie)
		numRows++
	}

	if numRows == 0 {
		return nil, errors.EntityNotFound{Entity: "movies"}
	}

	return movies, nil
}

func (dbstore *MovieStore) Create(ctx *gofr.Context, m *model.Movie) (*model.Movie, error) {
	stmt, err := ctx.DB().PrepareContext(ctx, `INSERT INTO MOVIES(NAME,GENRE,RATING,`+
		`RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED) VALUES(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return nil, errors.DB{Err: err}
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, m.Name, m.Genre, m.Rating, m.ReleaseDate, time.Now(), time.Now(), m.Plot, m.Released)

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	insertedID, _ := res.LastInsertId()
	movie, err := dbstore.GetByID(ctx, int(insertedID))

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (dbstore *MovieStore) Update(ctx *gofr.Context, m *model.Movie) (*model.Movie, error) {
	query, values := GenerateQuery(m)
	if values == nil {
		return nil, errors.Error("no fields to update")
	}

	stmt, err := ctx.DB().PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.DB{Err: err}
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, values...)

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return nil, errors.EntityNotFound{Entity: "id", ID: fmt.Sprint(m.ID)}
	}

	movie, _ := dbstore.GetByID(ctx, m.ID)

	return movie, nil
}
func (dbstore *MovieStore) Delete(ctx *gofr.Context, id int) error {
	stmt, err := ctx.DB().PrepareContext(ctx, "UPDATE MOVIES SET DELETEAT = ? WHERE ID = ? AND DELETEAT IS NULL")
	if err != nil {
		return errors.DB{Err: err}
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, time.Now(), id)

	if err != nil {
		return errors.DB{Err: err}
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.EntityNotFound{Entity: "id", ID: fmt.Sprint(id)}
	}

	return nil
}
