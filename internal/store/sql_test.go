package store

import (
	"context"
	"database/sql"
	"log"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/datastore"
	gofrerr "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shivam/Crud_Gofr/internal/models"
)

func newMock() (db *sql.DB, mock sqlmock.Sqlmock, store *MovieStore, ctx *gofr.Context) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Println(err)
	}

	store = New()
	ctx = gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})

	ctx.Context = context.Background()

	return
}

func Test_GETBYID(t *testing.T) {
	db, mock, MovieStore, ctx := newMock()

	defer db.Close()

	CurrTime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releaseDate", "updatedAt", "createdAt", "plot", "released"}).
		AddRow(1, "MazeRunner", "Action", 5, "2022-12-17", CurrTime.String(), CurrTime.String(), "PLOT!", true)

	testCases := []struct {
		desc   string
		id     int
		output *models.Movie
		query  interface{}
		err    error
	}{
		// success
		{
			desc: "Success",
			id:   1,
			output: &models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			},
			err: nil,
			query: mock.ExpectQuery(
				"select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where id = ? and deletedAt is null;").
				WithArgs(1).WillReturnRows(rows),
		},

		// error
		{
			desc:   "Fail",
			id:     1,
			output: nil,
			err:    gofrerr.Error("error while scanning the rows"),
			query: mock.ExpectQuery(
				"select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where id = ? and deletedAt is null;").
				WithArgs(1).WillReturnError(gofrerr.Error("error while scanning the rows")),
		},
	}

	for _, tc := range testCases {
		tc := tc

		output, err := MovieStore.GetByID(ctx, tc.id)
		if err != nil && (err.Error() != tc.err.Error()) {
			t.Errorf("Expected %v got %v", tc.err, err)
		}

		if !reflect.DeepEqual(output, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, output)
		}
	}
}

func Test_GETALL(t *testing.T) {
	db, mock, MovieStore, ctx := newMock()

	defer db.Close()

	CurrTime := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releaseDate", "updatedAt", "createdAt", "plot", "released"}).
		AddRow(1, "MazeRunner", "Action", 5, "2022-12-17", CurrTime.String(), CurrTime.String(), "PLOT!", true)

	testCases := []struct {
		desc   string
		id     int
		output []*models.Movie
		query  interface{}
		err    error
	}{
		// success
		{
			desc: "Success",
			id:   1,
			output: []*models.Movie{
				{
					ID:          1,
					Name:        "MazeRunner",
					Genre:       "Action",
					Rating:      5,
					ReleaseDate: "2022-12-17",
					UpdatedAt:   CurrTime.String(),
					CreatedAt:   CurrTime.String(),
					Plot:        "PLOT!",
					Released:    true,
				},
			},
			err: nil,
			query: mock.ExpectQuery(
				"select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where deletedAT is null;").
				WithArgs().WillReturnRows(rows),
		},

		// error
		{
			desc:   "Fail",
			id:     999999,
			output: nil,
			err:    gofrerr.Error("error while fetching the rows"),
			query: mock.ExpectQuery(
				"select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where deletedAT is null;").
				WithArgs().WillReturnError(gofrerr.Error("error while fetching the rows")),
		},

		// error
		{
			desc:   "Fail",
			id:     1,
			output: nil,
			err:    gofrerr.Error("error while scanning the rows"),
			query: mock.ExpectQuery(
				"select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where deletedAT is null;").
				WithArgs().
				WillReturnRows(sqlmock.NewRows([]string{"name", "genre", "rating", "releaseDate", "updatedAt", "createdAt", "plot", "released"}).
					AddRow("MazeRunner", "Action", 5, "2022-12-17", CurrTime.String(), CurrTime.String(), "PLOT!", true)),
		},
	}

	for _, tc := range testCases {
		tc := tc

		output, err := MovieStore.GetAll(ctx)
		if err != nil && (err.Error() != tc.err.Error()) {
			t.Errorf("Expected %v got %v", tc.err, err)
		}

		if !reflect.DeepEqual(output, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, output)
		}
	}
}

func Test_DeleteById(t *testing.T) {
	db, mock, MovieStore, ctx := newMock()
	defer db.Close()

	testcase := []struct {
		desc  string
		id    int
		err   error
		query interface{}
	}{
		{
			desc: "Success",
			id:   1,
			err:  nil,
			query: mock.ExpectExec(
				"update MOVIE set deletedAT = ? where id = ? AND deletedAT IS null;").
				WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1)),
		},

		{
			desc: "Fail",
			id:   1,
			err:  gofrerr.Error("failed to delete"),
			query: mock.ExpectExec(
				"update MOVIE set deletedAT = ? where id = ? AND deletedAT IS null;").
				WithArgs(sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 0)),
		},

		{
			desc: "Fail",
			id:   1,
			err:  gofrerr.Error("error while deleting"),
			query: mock.ExpectExec(
				"update MOVIE set deletedAT = ?  where id = ? AND deletedAT IS null;").
				WithArgs(time.Now(), 1).WillReturnError(gofrerr.Error("error while deleting")),
		},
	}

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			err := MovieStore.DeleteByID(ctx, tc.id)
			if err != nil && (err.Error() != tc.err.Error()) {
				t.Errorf("Expected %v got %v", tc.err, err)
			}
		})
	}
}

func Test_Create(t *testing.T) {
	db, mock, MovieStore, ctx := newMock()
	defer db.Close()

	testcase := []struct {
		desc   string
		input  *models.Movie
		output *models.Movie
		err    error
		query  interface{}
	}{
		{
			desc: "Success",
			input: &models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5.0,
				ReleaseDate: "2022-12-17",
				Plot:        "PLOT!",
				Released:    true,
			},
			output: &models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5.0,
				ReleaseDate: "2022-12-17",
				Plot:        "PLOT!",
				Released:    true,
			},
			err: nil,
			query: mock.
				ExpectExec("insert into MOVIE(id,name,genre,rating,releaseDate,plot,released) values(?,?,?,?,?,?,?);").
				WithArgs(1, "MazeRunner", "Action", 5.0, "2022-12-17", "PLOT!", true).
				WillReturnResult(sqlmock.NewResult(1, 1)),
		},

		{
			desc: "Fail",
			input: &models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				Plot:        "PLOT!",
				Released:    true,
			},
			output: nil,
			err:    gofrerr.Error("error while executing create query"),
			query: mock.
				ExpectExec("insert into MOVIE(id,name,genre,rating,releaseDate,plot,released) values(?,?,?,?,?,?,?);").
				WithArgs(1, "MazeRunner", "Action", 5, "2022-12-17", "PLOT!", true).
				WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(gofrerr.Error("error while executing create query")),
		},
	}

	for _, tc := range testcase {
		tc := tc

		out, err := MovieStore.Create(ctx, tc.input)
		if err != nil && (err.Error() != tc.err.Error()) {
			t.Errorf("Expected %v got %v", tc.err, err)
		}

		if !reflect.DeepEqual(out, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, out)
		}
	}
}

func TestUpdate(t *testing.T) {
	db, mock, MovieStore, ctx := newMock()
	rows := sqlmock.NewRows([]string{"id", "name", "genre", "rating", "releaseDate", "updatedAt", "createdAt", "plot", "released"}).
		AddRow(1, "MazeRunner", "Action", 5.0, "2022-12-17", "", "", "PLOT!", true)

	defer db.Close()

	testcase := []struct {
		desc   string
		id     int
		input  models.Movie
		output *models.Movie
		erout  error
		mock   interface{}
	}{

		{
			desc: "SUCCESS",
			id:   1,
			input: models.Movie{
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5.0,
				ReleaseDate: "2022-12-17",
				Plot:        "PLOT!",
			},
			output: &models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5.0,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   "",
				CreatedAt:   "",
				Plot:        "PLOT!",
				Released:    true,
			},
			erout: nil,
			mock: []interface{}{mock.ExpectExec(
				"UPDATE MOVIE SET  name = ?, genre = ?, rating = ?, releaseDate = ?, plot = ? where id = ? AND deletedAt IS NULL;").
				WithArgs("MazeRunner", "Action", 5.0, "2022-12-17", "PLOT!", 1).
				WillReturnResult(sqlmock.NewResult(0, 1)),
				mock.ExpectQuery(
					"select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where id = ? and deletedAt is null;").
					WithArgs(1).WillReturnRows(rows)},
		},
		{
			desc: "Failed Case",
			id:   1,
			input: models.Movie{
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5.0,
				ReleaseDate: "2022-12-17",
				Plot:        "PLOT!",
			},
			output: nil,
			erout:  gofrerr.Error("internal server error"),
			mock: []interface{}{mock.ExpectExec(
				"UPDATE MOVIE SET  name = ?, genre = ?, rating = ?, releaseDate = ?, plot = ? where id = ? AND deletedAt IS NULL;").
				WithArgs("MazeRunner", "Action", 5.0, "2022-12-17", "PLOT!", 1).
				WillReturnResult(sqlmock.NewResult(0, 1)),
				mock.ExpectQuery(
					"select id,name,genre,rating,releaseDate,updatedAt,createdAt,plot,released from MOVIE where id = ? and deletedAt is null;").
					WithArgs(1).WillReturnError(gofrerr.Error("error while scanning the rows"))},
		},
		{
			desc: "Failed Case",
			id:   1,
			input: models.Movie{
				Name: "MazeRunner",
			},
			output: nil,
			erout:  gofrerr.Error("error while executing the query for update"),
			mock: mock.ExpectExec("UPDATE MOVIE SET  name = ? where id = ? AND deletedAt IS NULL;").
				WithArgs("MazeRunner", 1).
				WillReturnError(gofrerr.Error("error while executing the query for update")),
		},
	}

	for _, tc := range testcase {
		tc := tc

		out, err := MovieStore.Update(ctx, tc.id, &tc.input)
		if err != nil && (err.Error() != tc.erout.Error()) {
			t.Errorf("Expected %v got %v", tc.erout, err)
		}

		if !reflect.DeepEqual(out, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, out)
		}
	}
}
