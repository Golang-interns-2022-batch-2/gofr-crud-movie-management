package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iamkakashi/movie-gofr/internal/model"
	"gopkg.in/guregu/null.v4"
)

func newMock(*testing.T) (db *sql.DB, mock sqlmock.Sqlmock, store *MovieStore, ctx *gofr.Context) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
	}

	store = New()
	ctx = gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
	ctx.Context = context.Background()

	return
}
func TestGetByID(t *testing.T) {
	db, mock, movieStore, ctx := newMock(t)

	defer db.Close()

	currenttime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releasedate", "updatedat", "createdat", "plot", "released"}).
		AddRow(4, "Deadpool", "Comedy, Action", 4.8, "2014-12-17", currenttime, currenttime, "This is a superhero movie", true)
	testcase := []struct {
		desc   string
		input  int
		output *model.Movie
		erout  error
		mock   interface{}
	}{
		{
			desc:  "Success",
			input: 4,
			output: &model.Movie{
				ID:          4,
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			erout: nil,
			mock: mock.ExpectQuery(
				"SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL AND ID = ?").
				WithArgs(4).WillReturnRows(rows),
		},
		{
			desc:   "Fail",
			input:  4,
			output: nil,
			erout:  errors.Error("sql: expected 8 destination arguments in Scan, not 9"),
			mock: mock.ExpectQuery(
				"SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL AND ID = ?").
				WithArgs(4).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releasedate", "updatedat", "createdat", "plot"}).
				AddRow(4, "Deadpool", "Comedy, Action", 4.8, "2014-12-17", currenttime, currenttime, "This is a superhero movie")),
		},
		{
			desc:   "Id not found",
			input:  4,
			output: nil,
			erout:  errors.EntityNotFound{Entity: "movie", ID: fmt.Sprint(4)},
			mock: mock.ExpectQuery(
				"SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL AND ID = ?").
				WithArgs(4).WillReturnError(errors.Error("sql: no rows in result set")),
		},
	}

	for _, tc := range testcase {
		out, err := movieStore.GetByID(ctx, tc.input)
		if err != nil && (err.Error() != tc.erout.Error()) {
			t.Errorf("Expected %v got %v", tc.erout, err)
		}

		if !reflect.DeepEqual(out, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, out)
		}
	}
}

func TestGet(t *testing.T) {
	db, mock, movieStore, ctx := newMock(t)

	defer db.Close()

	currenttime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releasedate", "updatedat", "createdat", "plot", "released"}).
		AddRow(4, "Deadpool", "Comedy, Action", 4.8, "2014-12-17", currenttime, currenttime, "This is a superhero movie", true)
	testcase := []struct {
		desc   string
		output []*model.Movie
		erout  error
		mock   interface{}
	}{
		{
			desc: "Success",
			output: []*model.Movie{{
				ID:          4,
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			},
			erout: nil,
			mock: mock.ExpectQuery("SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL").
				WithArgs().WillReturnRows(rows),
		},
		{
			desc:   "Scan Fail",
			output: nil,
			erout:  errors.Error("sql: expected 8 destination arguments in Scan, not 9"),
			mock: mock.
				ExpectQuery("SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL").
				WithArgs().
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releasedate", "updatedat", "createdat", "plot"}).
					AddRow(4, "Deadpool", "Comedy, Action", 4.8, "2014-12-17", currenttime, currenttime, "This is a superhero movie")),
		},
		{
			desc:   "Query Fail",
			output: nil,
			erout:  errors.Error("internal server error"),
			mock: mock.
				ExpectQuery("SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL").
				WithArgs().WillReturnError(errors.Error("internal server error")),
		},
		{
			desc:   "No row exist",
			output: nil,
			erout:  errors.EntityNotFound{Entity: "movies"},
			mock: mock.
				ExpectQuery("SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL").
				WithArgs().
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releasedate", "updatedat", "createdat", "plot", "released"})),
		},
	}

	for _, tc := range testcase {
		out, err := movieStore.Get(ctx)
		if err != nil && (err.Error() != tc.erout.Error()) {
			t.Errorf("Expected %v got %v", tc.erout, err)
		}

		if !reflect.DeepEqual(out, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, out)
		}
	}
}

func TestCreate(t *testing.T) {
	db, mock, movieStore, ctx := newMock(t)
	currenttime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releasedate", "updatedat", "createdat", "plot", "released"}).
		AddRow(4, "Deadpool", "Comedy, Action", 4.8, "2014-12-17", currenttime, currenttime, "This is a superhero movie", true)

	defer db.Close()

	testcase := []struct {
		desc   string
		input  *model.Movie
		output *model.Movie
		erout  error
		mock   interface{}
	}{
		{
			desc: "Success",
			input: &model.Movie{
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: &model.Movie{
				ID:          4,
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			erout: nil,
			mock: []interface{}{
				mock.ExpectPrepare("INSERT INTO MOVIES(NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED) VALUES(?,?,?,?,?,?,?,?)").
					ExpectExec().
					WithArgs("Deadpool", "Comedy, Action", 4.8, "2014-12-17", sqlmock.AnyArg(), sqlmock.AnyArg(), "This is a superhero movie", true).
					WillReturnResult(sqlmock.NewResult(4, 1)),
				mock.ExpectQuery(
					"SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL AND ID = ?").
					WithArgs(4).WillReturnRows(rows)},
		},
		{
			desc: "Prepare Failure",
			input: &model.Movie{
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: nil,
			erout:  errors.Error("error preparing create query"),
			mock: mock.ExpectPrepare("INSERT INTO MOVIES(NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED) VALUES(?,?,?,?,?,?,?,?)").
				WillReturnError(errors.Error("error preparing create query")),
		},
		{
			desc: "Execution Failure",
			input: &model.Movie{
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: nil,
			erout:  errors.Error("error executing create query"),
			mock: mock.ExpectPrepare(
				"INSERT INTO MOVIES(NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED) VALUES(?,?,?,?,?,?,?,?)").
				ExpectExec().
				WithArgs("Deadpool", "Comedy, Action", 4.8, "2014-12-17", sqlmock.AnyArg(), sqlmock.AnyArg(), "This is a superhero movie", true).
				WillReturnError(errors.Error("error executing create query")),
		},
		{
			desc: "Fail",
			input: &model.Movie{
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: nil,
			erout:  errors.Error("failed to create movie"),
			mock: []interface{}{
				mock.ExpectPrepare("INSERT INTO MOVIES(NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED) VALUES(?,?,?,?,?,?,?,?)").
					ExpectExec().
					WithArgs("Deadpool", "Comedy, Action", 4.8, "2014-12-17", sqlmock.AnyArg(), sqlmock.AnyArg(), "This is a superhero movie", true).
					WillReturnResult(sqlmock.NewResult(4, 1)),
				mock.ExpectQuery(
					"SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL AND ID = ?").
					WithArgs(4).WillReturnError(errors.Error("failed to create movie"))},
		},
	}

	for _, tc := range testcase {
		out, err := movieStore.Create(ctx, tc.input)
		if err != nil && (err.Error() != tc.erout.Error()) {
			t.Errorf("Expected %v got %v", tc.erout, err)
		}

		if !reflect.DeepEqual(out, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, out)
		}
	}
}

func TestDelete(t *testing.T) {
	db, mock, movieStore, ctx := newMock(t)
	defer db.Close()

	testcase := []struct {
		desc  string
		input int
		erout error
		mock  interface{}
	}{

		{
			desc:  "Success",
			input: 4,
			erout: nil,
			mock: mock.ExpectPrepare(
				"UPDATE MOVIES SET DELETEAT = ? WHERE ID = ? AND DELETEAT IS NULL").
				ExpectExec().WithArgs(sqlmock.AnyArg(), 4).WillReturnResult(sqlmock.NewResult(4, 1)),
		},
		{
			desc:  "Fail",
			input: 4,
			erout: errors.EntityNotFound{Entity: "id", ID: "4"},
			mock: mock.ExpectPrepare(
				"UPDATE MOVIES SET DELETEAT = ? WHERE ID = ? AND DELETEAT IS NULL").
				ExpectExec().WithArgs(sqlmock.AnyArg(), 4).WillReturnResult(sqlmock.NewResult(0, 0)),
		},

		{
			desc:  "Execution Failure",
			input: 4,
			erout: errors.Error("error executing delete query"),
			mock: mock.ExpectPrepare(
				"UPDATE MOVIES SET DELETEAT = ? WHERE ID = ? AND DELETEAT IS NULL").
				ExpectExec().WithArgs(sqlmock.AnyArg(), 4).WillReturnError(errors.Error("error executing delete query")),
		},
		{
			desc:  "Prepare Failure",
			input: 4,
			erout: errors.Error("error preparing delete query"),
			mock: mock.ExpectPrepare(
				"UPDATE MOVIES SET DELETEAT = ? WHERE ID = ? AND DELETEAT IS NULL").
				WillReturnError(errors.Error("error preparing delete query")),
		},
	}

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			err := movieStore.Delete(ctx, tc.input)
			if err != nil && (err.Error() != tc.erout.Error()) {
				t.Errorf("Expected %v got %v", tc.erout, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db, mock, movieStore, ctx := newMock(t)
	currenttime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releasedate", "updatedat", "createdat", "plot", "released"}).
		AddRow(4, "Deadpool", "Comedy, Action", 4.8, "2014-12-17", currenttime, currenttime, "This is a superhero movie", true)

	defer db.Close()

	testcase := []struct {
		desc   string
		input  *model.Movie
		output *model.Movie
		erout  error
		mock   interface{}
	}{
		{
			desc: "Failed Case",
			input: &model.Movie{
				ID:    4,
				Name:  null.NewString("Deadpool", true),
				Genre: null.NewString("Comedy, Action", true),
			},
			output: nil,
			erout:  errors.EntityNotFound{Entity: "id", ID: "4"},
			mock: mock.ExpectPrepare("UPDATE MOVIES SET NAME = ?, GENRE = ?, UPDATEDAT=NOW() WHERE ID = ? AND DELETEAT IS NULL").
				ExpectExec().
				WithArgs("Deadpool", "Comedy, Action", 4).
				WillReturnResult(sqlmock.NewResult(0, 0)),
		},
		{
			desc: "SUCCESS",
			input: &model.Movie{
				ID:   4,
				Name: null.NewString("Deadpool", true),
			},
			output: &model.Movie{
				ID:          4,
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			erout: nil,
			mock: []interface{}{mock.ExpectPrepare("UPDATE MOVIES SET NAME = ?, UPDATEDAT=NOW() WHERE ID = ? AND DELETEAT IS NULL").
				ExpectExec().
				WithArgs("Deadpool", 4).
				WillReturnResult(sqlmock.NewResult(0, 1)), mock.ExpectQuery(
				"SELECT ID,NAME,GENRE,RATING,RELEASEDATE,UPDATEDAT,CREATEDAT,PLOT,RELEASED FROM MOVIES WHERE DELETEAT IS NULL AND ID = ?",
			).WithArgs(4).WillReturnRows(rows)},
		},
		{
			desc: "No fields to update",
			input: &model.Movie{
				ID: 4,
			},
			output: nil,
			erout:  errors.Error("no fields to update"),
		},
		{
			desc: "Prepare Fail",
			input: &model.Movie{
				ID:    4,
				Name:  null.NewString("Deadpool", true),
				Genre: null.NewString("Comedy, Action", true),
			},
			output: nil,
			erout:  errors.Error("error preparing update query"),
			mock: mock.ExpectPrepare("UPDATE MOVIES SET NAME = ?, GENRE = ?, UPDATEDAT=NOW() WHERE ID = ? AND DELETEAT IS NULL").
				WillReturnError(errors.Error("error preparing update query")),
		},
		{
			desc: "Execution Fail",
			input: &model.Movie{
				ID:          4,
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: nil,
			erout:  errors.Error("error executing update query"),
			mock: mock.ExpectPrepare(
				`UPDATE MOVIES SET NAME = ?, GENRE = ?, PLOT = ?, RATING = ?, RELEASEDATE = ?, RELEASED = ?, UPDATEDAT=NOW()
				WHERE ID = ? AND DELETEAT IS NULL`,
			).
				ExpectExec().
				WithArgs("Deadpool", "Comedy, Action", "This is a superhero movie", 4.8, "2014-12-17", true, 4).
				WillReturnError(errors.Error("error executing update query")),
		},
	}

	for _, tc := range testcase {
		out, err := movieStore.Update(ctx, tc.input)
		if err != nil && (err.Error() != tc.erout.Error()) {
			t.Errorf("Expected %v got %v", tc.erout, err)
		}

		if !reflect.DeepEqual(out, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, out)
		}
	}
}
