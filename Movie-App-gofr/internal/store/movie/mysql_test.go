package movie

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/model"
)

func TestCreateMovie(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf(err.Error())
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})

	ctx.Context = context.Background()
	defer db.Close()

	query := "insert into movie_details(name,genre,rating,plot,released,releaseDate) values(?,?,?,?,?,?); "

	tcs := []struct {
		Desc        string
		ID          int
		Name        string
		Genre       string
		Rating      float64
		Plot        string
		released    bool
		releaseDate string
		mockQ       interface{}
		expectedErr error
	}{
		{
			Desc:        "Success",
			ID:          1,
			Name:        "Silicon Valley",
			Genre:       "comedy",
			Rating:      4.5,
			Plot:        "Richard",
			released:    true,
			releaseDate: "2014-12-17",
			mockQ: mock.
				ExpectExec(query).
				WithArgs("Silicon Valley", "comedy", 4.5, "Richard", true, "2014-12-17").
				WillReturnResult(sqlmock.NewResult(0, 1)),
			expectedErr: nil,
		},
		{
			Desc:        "Prepare Error Tc",
			ID:          1,
			Name:        "Silicon Valley",
			Genre:       "comedy",
			Rating:      4.5,
			Plot:        "Richard",
			released:    true,
			releaseDate: "2014-12-17",
			mockQ: mock.
				ExpectExec(query).
				WithArgs("Silicon Valley", "comedy", 4.5, "Richard", true, "2014-12-17").
				WillReturnError(errors.Error("Internal Server Error")),
			expectedErr: errors.Error("Internal Server Error"),
		},
		{
			Desc:        "Failure",
			ID:          1,
			Name:        "Silicon Valley",
			Genre:       "comedy",
			Rating:      4.5,
			Plot:        "Richard",
			released:    true,
			releaseDate: "2014-12-17",
			mockQ: mock.
				ExpectExec(query).
				WithArgs("Silicon Valley", "comedy", 4.5, "Richard", true, "2014-12-17").
				WillReturnError(nil),
			expectedErr: nil,
		},
	}
	handler := NewDBHandler()

	for _, tt := range tcs {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			tmpObj := model.MovieModel{
				ID: tt.ID, Name: tt.Name, Genre: tt.Genre, Rating: tt.Rating,
				Plot: tt.Plot, ReleaseDate: tt.releaseDate, Released: tt.released,
			}
			mObj, err := handler.CreateMovie(ctx, &tmpObj)

			if mObj != nil && err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Got : %v, Want : %v", err.Error(), tt.expectedErr.Error())
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf(err.Error())
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})

	ctx.Context = context.Background()

	query := "select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from movie_details where id = ? and deletedAt is null;"

	tcs := []struct {
		Desc        string
		ID          int
		Name        string
		Genre       string
		Rating      float64
		Plot        string
		released    bool
		releaseDate string
		updatedAt   string
		createdAt   string
		mockQ       interface{}
		expectedErr error
	}{
		{
			Desc:        "Success",
			ID:          1,
			Name:        "Silicon Valley",
			Genre:       "comedy",
			Rating:      4.5,
			Plot:        "Richard",
			released:    true,
			releaseDate: "2014-12-17",
			mockQ: mock.ExpectQuery(query).
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releaseDate", "updatedAt", "createdAt", "plot", "released"}).
					AddRow(1, "Silicon Valley", "comedy", 4.5, "2014-12-17", "2014-12-17", "2014-12-17", "Richard", true)),
			expectedErr: nil,
		},
		{
			Desc:        "Success",
			ID:          -1,
			Name:        "Silicon Valley",
			Genre:       "comedy",
			Rating:      4.5,
			Plot:        "Richard",
			released:    true,
			releaseDate: "2014-12-17",
			mockQ: mock.ExpectPrepare(query).
				ExpectQuery().
				WithArgs(-1).
				WillReturnError(sql.ErrNoRows),
			expectedErr: sql.ErrNoRows,
		},
		{
			Desc:        "Success",
			ID:          1000,
			Name:        "Silicon Valley",
			Genre:       "comedy",
			Rating:      4.5,
			Plot:        "Richard",
			released:    true,
			releaseDate: "2014-12-17",
			mockQ: mock.ExpectPrepare(query).
				ExpectQuery().
				WithArgs(1000).
				WillReturnError(sql.ErrNoRows),
			expectedErr: sql.ErrNoRows,
		},
		{
			Desc:        "Success",
			ID:          1,
			Name:        "Silicon Valley",
			Genre:       "comedy",
			Rating:      4.5,
			Plot:        "Richard",
			released:    true,
			releaseDate: "2014-12-17",
			mockQ: mock.ExpectQuery(query).
				WithArgs(1).WillReturnError(sql.ErrNoRows),
			expectedErr: sql.ErrNoRows,
		},
	}
	handler := NewDBHandler()

	for _, tt := range tcs {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			mObj, err := handler.GetByID(ctx, tt.ID)
			if mObj != nil && err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Got : %v, Want : %v", err.Error(), tt.expectedErr.Error())
			}
		})
	}
}

func TestUpdateByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf(err.Error())
	}

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})

	ctx.Context = context.Background()

	query := "update movie_details set rating=?,plot=?,releaseDate=?,updatedAt=? where deletedAt is null and id=?;"

	tcs := []struct {
		Desc        string
		ID          int
		Name        string
		Genre       string
		Rating      float64
		Plot        string
		released    bool
		releaseDate string
		updatedAt   string
		createdAt   string
		UpdatedAt   string
		mockQ       interface{}
		expectedErr error
	}{

		{
			Desc:        "exec success",
			ID:          2,
			Name:        "Silicon Valley",
			Genre:       "comedy",
			Rating:      4.5,
			Plot:        "Richard",
			released:    true,
			releaseDate: "2014-12-17",
			updatedAt:   "2014-12-17",
			mockQ: mock.ExpectExec(query).
				WithArgs(4.5, "Richard", sqlmock.AnyArg(), sqlmock.AnyArg(), 2).
				WillReturnError(sql.ErrNoRows),
			expectedErr: sql.ErrNoRows,
		},
	}

	handler := NewDBHandler()

	for _, tt := range tcs {
		tt := tt

		t.Run(tt.Desc, func(t *testing.T) {
			tmpObj := model.MovieModel{
				ID:          tt.ID,
				Name:        tt.Name,
				Genre:       tt.Genre,
				Rating:      tt.Rating,
				Plot:        tt.Plot,
				Released:    tt.released,
				ReleaseDate: tt.releaseDate,
				UpdatedAt:   tt.updatedAt}

			mObj, err := handler.UpdateByID(ctx, &tmpObj)

			if mObj != nil && err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Got : %v, Want : %v", err.Error(), tt.expectedErr.Error())
			}
		})
	}
}

func TestDeleteByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})

	ctx.Context = context.Background()

	if err != nil {
		t.Errorf(err.Error())
	}

	query := "update movie_details set deletedAt = ? where id = ? and deletedAt is null;"
	format := "2006-01-02 15:04:05"
	tcs := []struct {
		Desc        string
		ID          int
		Name        string
		Genre       string
		Rating      float64
		Plot        string
		released    bool
		releaseDate string
		updatedAt   string
		createdAt   string
		UpdatedAt   string
		mockQ       interface{}
		mockTime    time.Time
		expectedErr error
	}{

		{
			Desc:     "Success",
			ID:       1,
			Name:     "Silicon Valley",
			Genre:    "comedy",
			Rating:   4.5,
			Plot:     "Richard",
			released: true,

			mockQ: mock.ExpectExec(query).
				WithArgs(time.Now().Format(format), 1).
				WillReturnError(sql.ErrNoRows),
			expectedErr: sql.ErrNoRows,
		},
	}

	handler := NewDBHandler()

	for _, tt := range tcs {
		tt := tt
		err := handler.DeleteByID(ctx, tt.ID)
		t.Run(tt.Desc, func(t *testing.T) {
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Got : %v and Expected error: %v", err.Error(), tt.expectedErr.Error())
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})

	ctx.Context = context.Background()
	if err != nil {
		t.Errorf(err.Error())
	}

	query := "select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from movie_details where deletedAt is null;"

	tcs := []struct {
		Desc        string
		expectedErr error
		mockQ       interface{}
	}{

		{
			Desc:        "Exec success tc",
			expectedErr: nil,
			mockQ: mock.ExpectQuery(query).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releaseDate", "updatedAt", "createdAt", "plot", "released"}).
					AddRow(1, "Richard", "comedy", 4.5, "2022-03-06", "2022-03-06 11:23:54", "2022-03-06 11:23:54", "Richard", true)).
				WillReturnError(nil),
		},
		{
			Desc:        "Exec failure tc",
			expectedErr: errors.Error("Couldn't execute query"),
			mockQ:       mock.ExpectQuery(query).WillReturnError(errors.Error("Couldn't execute query")),
		},
		{
			Desc:        "Exec success tc",
			expectedErr: nil,
			mockQ: mock.
				ExpectQuery(query).
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releaseDate", "updatedAt", "createdAt", "plot", "released"}).
					AddRow(1, "Richard", "comedy", 4.5, "2022-03-06", "2022-03-06 11:23:54", "2022-03-06 11:23:54", "Richard", true)).
				WillReturnError(nil),
		},
	}

	handler := NewDBHandler()

	for _, tt := range tcs {
		tt := tt

		_, err := handler.GetAll(ctx)

		t.Run(tt.Desc, func(t *testing.T) {
			if err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Got : %v and Expected error: %v", err.Error(), tt.expectedErr.Error())
			}
		})
	}
}
