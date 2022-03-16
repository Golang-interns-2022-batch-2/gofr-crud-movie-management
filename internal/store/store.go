package store

import (
	"strings"
	"time"

	gofrerr "developer.zopsmart.com/go/gofr/pkg/errors"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/shivam/Crud_Gofr/internal/models"
)

type MovieStore struct{}

func New() *MovieStore {
	return &MovieStore{}
}

func QueryUpdate(movie *models.Movie) (query string, values []interface{}) {
	if movie.Name != "" {
		query += " name = ?, "

		values = append(values, movie.Name)
	}

	if movie.Genre != "" {
		query += " genre = ?, "

		values = append(values, movie.Genre)
	}

	if movie.Rating > 0 {
		query += " rating = ?, "

		values = append(values, movie.Rating)
	}

	if movie.ReleaseDate != "" {
		query += " releaseDate = ?, "

		values = append(values, movie.ReleaseDate)
	}

	if movie.Plot != "" {
		query += " plot = ?, "

		values = append(values, movie.Plot)
	}

	query = strings.TrimSuffix(query, ", ")
	query += " where id = ?"

	values = append(values, movie.ID)

	return query, values
}

func (db *MovieStore) Update(ctx *gofr.Context, id int, updatedDetails *models.Movie) (*models.Movie, error) {
	updatedDetails.ID = id
	query, values := QueryUpdate(updatedDetails)
	updateQuery := "UPDATE MOVIE SET "
	updateQuery = updateQuery + query + " AND deletedAt IS NULL;"

	_, err := ctx.DB().ExecContext(ctx, updateQuery, values...)

	if err != nil {
		return nil, gofrerr.Error("error while executing the query for update")
	}

	NewDetails, er := db.GetByID(ctx, id)

	if er == nil {
		return NewDetails, nil
	}

	return nil, gofrerr.Error("internal server error")
}
func (db *MovieStore) GetByID(ctx *gofr.Context, id int) (*models.Movie, error) {
	FoundRow := models.Movie{}

	que := "select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where id = ? and deletedAt is null;"

	row := ctx.DB().QueryRowContext(ctx, que, id)
	err := row.Scan(&FoundRow.ID, &FoundRow.Name, &FoundRow.Genre,
		&FoundRow.Rating, &FoundRow.ReleaseDate, &FoundRow.UpdatedAt,
		&FoundRow.CreatedAt, &FoundRow.Plot, &FoundRow.Released)

	if err != nil {
		return nil, err
	}

	return &FoundRow, nil
}

func (db *MovieStore) DeleteByID(ctx *gofr.Context, id int) error {
	que := "update MOVIE set deletedAT = ? where id = ? AND deletedAT IS null;"
	res, err := ctx.DB().ExecContext(ctx, que, time.Now(), id)

	if err != nil {
		return gofrerr.Error("error while deleting")
	}

	number, _ := res.RowsAffected()
	if number == 0 {
		return gofrerr.Error("failed to delete")
	}

	return nil
}

func (db MovieStore) GetAll(ctx *gofr.Context) ([]*models.Movie, error) {
	que := "select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where deletedAT is null;"
	rows, err := ctx.DB().QueryContext(ctx, que)

	if err != nil {
		return nil, gofrerr.Error("error while fetching the rows")
	}

	Entries := []*models.Movie{}

	for rows.Next() {
		FoundRow := models.Movie{}
		err := rows.Scan(&FoundRow.ID, &FoundRow.Name, &FoundRow.Genre,
			&FoundRow.Rating, &FoundRow.ReleaseDate, &FoundRow.UpdatedAt,
			&FoundRow.CreatedAt, &FoundRow.Plot, &FoundRow.Released)

		Entries = append(Entries, &FoundRow)

		if err != nil {
			return nil, gofrerr.Error("error while scanning the rows")
		}
	}

	return Entries, nil
}

func (db *MovieStore) Create(ctx *gofr.Context, a *models.Movie) (*models.Movie, error) {
	que := "insert into MOVIE(id,name,genre,rating,releaseDate,plot,released) values(?,?,?,?,?,?,?); "

	_, err1 := ctx.DB().ExecContext(ctx, que, a.ID, a.Name, a.Genre, a.Rating, a.ReleaseDate, a.Plot, a.Released)
	if err1 != nil {
		return nil, gofrerr.Error("error while executing create query")
	}

	return a, nil
}
