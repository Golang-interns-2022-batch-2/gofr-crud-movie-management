package movie

import (
	"fmt"
	"log"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/models"

	"developer.zopsmart.com/go/gofr/pkg/errors"
)

type Store struct {
}

func New() *Store {
	return &Store{}
}

func (s *Store) GetByID(ctx *gofr.Context, id int) (*models.Movie, error) {
	var mov models.Movie

	query := "select id , name ,genre ,rating , releaseDate ,updatedAt,createdAt,plot,released from movie where deletedAt IS NULL and id=?"
	row := ctx.DB().QueryRowContext(ctx, query, id)

	err := row.Scan(&mov.ID, &mov.Name, &mov.Genre, &mov.Rating, &mov.ReleaseDate, &mov.UpdatedAt, &mov.CreatedAt, &mov.Plot, &mov.Released)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, errors.EntityNotFound{Entity: "movie", ID: fmt.Sprint(id)}
	}

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return &mov, nil
}
func (s *Store) Delete(ctx *gofr.Context, id int) error {
	query := "update  movie set deletedAt=Now() where id=?"
	result, err := ctx.DB().ExecContext(ctx, query, id)

	if err != nil {
		return errors.DB{Err: err}
	}

	res, _ := result.RowsAffected()
	if res == 0 {
		return errors.EntityNotFound{Entity: "id", ID: fmt.Sprint(id)}
	}

	return nil
}

func (s *Store) Update(ctx *gofr.Context, mov *models.Movie) (*models.Movie, error) {
	query := "update movie set name = ?," +
		" genre=?, rating=? ,releaseDate=?," +
		"updatedAt=?, createdAt=? ,plot=?," +
		" released=? WHERE deletedAt IS NULL" +
		" and id = ?"

	_, err := ctx.DB().ExecContext(ctx, query, &mov.Name,
		&mov.Genre, &mov.Rating,
		&mov.ReleaseDate, &mov.UpdatedAt,
		&mov.CreatedAt, &mov.Plot,
		&mov.Released, &mov.ID)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	return s.GetByID(ctx, mov.ID)
}

func (s *Store) Create(ctx *gofr.Context, mov *models.Movie) (*models.Movie, error) {
	res, err := ctx.DB().ExecContext(ctx, "insert into movie "+
		"(name,genre,rating,releaseDate,"+
		"updatedAt,createdAt,plot,released) "+
		"VALUES (?, ?, ?,?,?,?,?,?)", mov.Name,
		mov.Genre, mov.Rating, mov.ReleaseDate,
		time.Now(), time.Now(), mov.Plot, mov.Released)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	id, _ := res.LastInsertId()

	return s.GetByID(ctx, int(id))
}

func (s *Store) GetAll(ctx *gofr.Context) ([]*models.Movie, error) {
	var mov []*models.Movie

	result, err := ctx.DB().QueryContext(ctx, "SELECT id,name,"+
		"genre,rating,releaseDate,updatedAt,"+
		"createdAt,plot,released FROM movie "+
		"where deletedAt IS NULL")
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer func() {
		err = result.Err()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}()

	defer func() {
		err = result.Close()
		if err != nil {
			log.Printf("error: %v", err)
		}
	}()

	for result.Next() {
		var c models.Movie

		err = result.Scan(&c.ID, &c.Name, &c.Genre, &c.Rating, &c.ReleaseDate, &c.UpdatedAt, &c.CreatedAt, &c.Plot, &c.Released)
		if err != nil {
			return nil, errors.DB{Err: err}
		}

		mov = append(mov, &c)
	}

	return mov, nil
}
