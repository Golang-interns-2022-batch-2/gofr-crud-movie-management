package moviestore

import (
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"golangprog/gofr-test/Internal/models"
	"time"
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

func (s *Store) GetByID(ctx *gofr.Context, id int) (*models.Movie, error) {

	movieObj := models.Movie{}
	row := ctx.DB().QueryRowContext(ctx, "Select id, name, genre, rating,"+
		"	releaseDate, updatedAt, createdAt, plot, released FROM movie WHERE id = ? "+
		"and deletedAt IS NULL", id)

	err := row.Scan(&movieObj.ID, &movieObj.Name, &movieObj.Genre,
		&movieObj.Rating, &movieObj.ReleaseDate, &movieObj.UpdatedAt, &movieObj.CreatedAt,
		&movieObj.Plot, &movieObj.Released)
	if err != nil {
		return nil, err
	}

	return &movieObj, nil
}

func (s *Store) GetAll(ctx *gofr.Context) ([]*models.Movie, error) {
	movieObjs := []*models.Movie{}

	rows, err := ctx.DB().Query("Select id, name, genre, rating, releaseDate, " +
		"updatedAt, createdAt, plot, released FROM movie WHERE deletedAt IS NULL")
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	if err != nil {
		return nil, err
	}
	numRows := 0
	defer rows.Close()

	for rows.Next() {
		movieObj := models.Movie{}
		if err := rows.Scan(&movieObj.ID, &movieObj.Name, &movieObj.Genre,
			&movieObj.Rating, &movieObj.ReleaseDate, &movieObj.UpdatedAt,
			&movieObj.CreatedAt, &movieObj.Plot, &movieObj.Released); err != nil {
			return nil, err
		}

		movieObjs = append(movieObjs, &movieObj)
		numRows++
	}
	if numRows == 0 {
		return nil, err
	}
	return movieObjs, nil
}

func (s *Store) Delete(ctx *gofr.Context, id int) error {
	row, err := ctx.DB().Exec("update movie set deletedAt = Now() Where id = ?", id)
	if err != nil {
		return err
	}

	n, _ := row.RowsAffected()

	if n == 0 {
		return err
	}

	return nil
}

func (s *Store) Create(ctx *gofr.Context, movieObj *models.Movie) (*models.Movie, error) {
	query := `"CREATE TABLE IF NOT EXISTS movie(id int auto_increment primary key, 
		name varchar(25), genre varchar(25), rating float, releaseDate datetime, 
		updatedAt TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE now(), 
		createdAt TIMESTAMP NOT NULL DEFAULT NOW(), plot varchar(255), released bool)"`

	_, _ = ctx.DB().ExecContext(ctx, query)

	res, err := ctx.DB().Exec("INSERT into movie (name, genre, rating, "+
		"releaseDate, createdAt, updatedAt, plot, released) VALUES (?,?,?,?,?,?,?,?)", movieObj.Name,
		movieObj.Genre, movieObj.Rating, movieObj.ReleaseDate, time.Now(), time.Now(),
		movieObj.Plot, movieObj.Released)

	if err != nil {
		return nil, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return movieObj, nil
}

func (s *Store) Update(ctx *gofr.Context, movieObj *models.Movie) (*models.Movie, error) {
	row, err := ctx.DB().Exec("Update movie SET name=?, genre=?, rating=?, "+
		"releaseDate=?, updatedAt=?, createdAt=?, plot=?, released=? WHERE id = ? "+
		"and deletedAt IS NULL", movieObj.Name, movieObj.Genre, movieObj.Rating,
		movieObj.ReleaseDate, movieObj.CreatedAt, time.Now(), movieObj.Plot,
		movieObj.Released, movieObj.ID)
	if err != nil {
		return nil, err
	}

	n, _ := row.RowsAffected()

	if n == 0 {
		return nil, err
	}

	return movieObj, nil
}
