package movie

import (
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/models"
	"reflect"
	"testing"
	"time"

	"context"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	date, _ := time.Parse(time.RFC3339, "2014-12-17")
	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"ID", "Name",
		"Genre", "Rating", "ReleaseDate", "UpdatedAt",
		"CreatedAt", "Plot", "Released"}).AddRow(3,
		"Silicon valley", "Comedy", 4.5, date,
		dateTime, dateTime, "Richard, a programmer,"+
			" creates an app called the Pied Piper"+
			" and tries to get investors for it."+
			" Meanwhile, five other programmers "+
			"struggle to make their mark in Silicon Valley.", true)

	testCases := []struct {
		id            int
		mov           *models.Movie
		mockQuery     interface{}
		expectedError error
	}{
		{id: 3,
			mov: &models.Movie{ID: 3, Name: "Silicon valley",
				Genre: "Comedy", Rating: 4.5, ReleaseDate: date,
				UpdatedAt: dateTime, CreatedAt: dateTime,
				Plot: "Richard, a programmer, creates an app " +
					"called the Pied Piper and tries to get " +
					"investors for it. Meanwhile, five other" +
					" programmers struggle to make their " +
					"mark in Silicon Valley.", Released: true},
			mockQuery: mock.ExpectQuery("select id , " +
				"name ,genre ,rating , releaseDate ," +
				"updatedAt,createdAt,plot,released " +
				"from movie where deletedAt IS " +
				"NULL and id=?").WithArgs(3).WillReturnRows(rows),
			expectedError: nil},
		{
			id:  4,
			mov: nil,
			mockQuery: mock.ExpectQuery(
				"select id , name ,genre ," +
					"rating , releaseDate ,updatedAt," +
					"createdAt,plot,released from " +
					"movie where deletedAt IS NULL " +
					"and id=?").WithArgs(4).WillReturnError(errors.Error("new")),

			expectedError: errors.DB{Err: errors.Error("new")},
		},
		{
			id:  7,
			mov: nil,
			mockQuery: mock.ExpectQuery("select id ," +
				" name ,genre ,rating , releaseDate ," +
				"updatedAt,createdAt,plot,released " +
				"from movie where deletedAt IS" +
				" NULL and id=?").WithArgs(7).
				WillReturnError(errors.Error("sql: no rows in result set")),
			expectedError: errors.EntityNotFound{Entity: "movie", ID: fmt.Sprint(7)}},
	}

	for _, testCase := range testCases {
		s := New()

		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()

		mov, err := s.GetByID(ctx, testCase.id)

		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected error %v got %v", testCase.expectedError, err)
		}

		if !reflect.DeepEqual(mov, testCase.mov) {
			t.Errorf("expected movie %v"+
				" got %v ", &models.Movie{ID: 3,
				Name: "Silicon valley", Genre: "Comedy",
				Rating: 4.5, ReleaseDate: date, UpdatedAt: dateTime,
				CreatedAt: dateTime, Plot: "Richard, a programmer, " +
					"creates an app called the Pied Piper and " +
					"tries to get investors for it. Meanwhile," +
					" five other programmers struggle to m" +
					"ake their mark in Silicon Valley.", Released: true}, mov)
		}
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	testCases := []struct {
		desc      string
		id        int
		expecterr error
		mockCall  *sqlmock.ExpectedExec
	}{
		{desc: "Delete success",
			id:        3,
			expecterr: nil,
			mockCall:  mock.ExpectExec("update  movie set deletedAt=Now() where id=?").WithArgs(3).WillReturnResult(sqlmock.NewResult(1, 1))},

		{
			desc:      "Negative ID ",
			id:        -1,
			expecterr: errors.DB{Err: errors.Error("new")},
			mockCall:  mock.ExpectExec("update  movie set deletedAt=Now() where id=?").WithArgs(-1).WillReturnError(errors.Error("new")),
		},

		{
			desc:      "Delete failed",
			id:        7,
			expecterr: errors.DB{Err: errors.Error("new")},
			mockCall: mock.ExpectExec("update  movie set deletedAt=Now() where id=?").WithArgs(7).
				WillReturnError(errors.Error("new")),
		},
	}

	for _, tc := range testCases {
		s := New()
		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()

		err := s.Delete(ctx, tc.id)

		if !reflect.DeepEqual(err, tc.expecterr) {
			t.Errorf("expected error %v got %v", tc.expecterr, err)
		}
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"ID", "Name", "Genre",
		"Rating", "ReleaseDate", "UpdatedAt",
		"CreatedAt", "Plot", "Released"}).
		AddRow(3, "Silicon valley",
			"Comedy", 4.5, date, dateTime,
			dateTime, "Richard, a programmer, "+
				"creates an app called the Pied Piper "+
				"and tries to get investors for it. Meanwhile,"+
				" five other programmers struggle to make "+
				"their mark in Silicon Valley.", true)

	tests := []struct {
		desc      string
		expecterr error
		inputMov  *models.Movie
		mockCall  interface{}
	}{
		{
			desc:      "update succes",
			expecterr: nil,
			inputMov: &models.Movie{ID: 3, Name: "Silicon valley",
				Genre: "Comedy", Rating: 4.5, ReleaseDate: date,
				UpdatedAt: dateTime, CreatedAt: dateTime,
				Plot: "Richard, a programmer, creates an " +
					"app called the Pied Piper and tries " +
					"to get investors for it. Meanwhile," +
					" five other programmers struggle to" +
					" make their mark in Silicon Valley.", Released: true},
			mockCall: []interface{}{mock.ExpectExec("update movie "+
				"set name = ?, genre=?, rating=? ,releaseDate=?,"+
				"updatedAt=?, createdAt=? ,plot=?, "+
				"released=? WHERE deletedAt IS "+
				"NULL and id = ?").WithArgs("Silicon valley",
				"Comedy", 4.5, date, dateTime, dateTime,
				"Richard, a programmer, creates an app called the "+
					"Pied Piper and tries to get investors for it."+
					" Meanwhile, five other programmers struggle "+
					"to make their mark in Silicon Valley.",
				true, 3).WillReturnResult(sqlmock.NewResult(1, 1)),
				mock.ExpectQuery("select id , name " +
					",genre ,rating ," +
					" releaseDate ,updatedAt,createdAt," +
					"plot,released from movie where" +
					" deletedAt IS NULL and id=?").
					WithArgs(3).WillReturnRows(rows),
			},
		},
		{
			desc:      "update fail",
			expecterr: errors.DB{Err: errors.Error("new")},
			inputMov: &models.Movie{ID: -4, Name: "",
				Genre: "", Rating: 0.0,
				ReleaseDate: time.Time{}, UpdatedAt: time.Time{},
				CreatedAt: time.Time{}, Plot: "",
				Released: false},
			mockCall: mock.ExpectExec("update movie "+
				"set name = ?, genre=?, rating=? ,"+
				"releaseDate=?,updatedAt=?, createdAt=? ,"+
				"plot=?, released=? WHERE deletedAt"+
				" IS NULL and id = ?").WithArgs("", "", 0.0,
				time.Time{}, sqlmock.AnyArg(), sqlmock.AnyArg(),
				"", false, -4).WillReturnError(errors.Error("new")),
		},
	}
	for _, tc := range tests {
		s := New()

		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()

		_, err := s.Update(ctx, tc.inputMov)

		if !reflect.DeepEqual(err, tc.expecterr) {
			t.Errorf("Expected: %v, Got: %v", tc.expecterr, err)
		}
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Now()
	rows := sqlmock.NewRows([]string{"ID", "Name", "Genre",
		"Rating", "ReleaseDate", "UpdatedAt",
		"CreatedAt", "Plot", "Released"}).
		AddRow(3, "Silicon valley",
			"Comedy", 4.5, date, dateTime, dateTime,
			"Richard, a programmer, creates an "+
				"app called the Pied Piper and tries "+
				"to get investors for it. Meanwhile,"+
				" five other programmers struggle "+
				"to make their mark in Silicon Valley.", true)

	tests := []struct {
		desc      string
		expecterr error
		inputMov  *models.Movie
		mockCall  interface{}
	}{
		{
			desc:      "Create Success",
			expecterr: nil,
			inputMov: &models.Movie{Name: "Silicon valley",
				Genre: "Comedy", Rating: 4.5,
				ReleaseDate: date, UpdatedAt: dateTime,
				CreatedAt: dateTime, Plot: "Richard, a programmer, " +
					"creates an app called the Pied " +
					"Piper and tries to get " +
					"investors for it. Meanwhile, " +
					"five other programmers struggle " +
					"to make their mark in Silicon" +
					" Valley.", Released: true},
			mockCall: []interface{}{mock.ExpectExec("insert into movie "+
				"(name,genre,rating,releaseDate,updatedAt,"+
				"createdAt,plot,released) VALUES (?, ?, ?,"+
				"?,?,?,?,?)").WithArgs("Silicon valley",
				"Comedy", 4.5, date, sqlmock.AnyArg(),
				sqlmock.AnyArg(), "Richard, a programmer, creates "+
					"an app called the Pied Piper and tries to"+
					" get investors for it. Meanwhile, five other"+
					" programmers struggle to make their "+
					"mark in Silicon Valley.", true).
				WillReturnResult(sqlmock.NewResult(3, 1)),
				mock.ExpectQuery("select id , name ,genre" +
					" ,rating , releaseDate ,updatedAt,createdAt,plot," +
					"released from movie where deletedAt" +
					" IS NULL and id=?").WithArgs(3).WillReturnRows(rows),
			},
		},

		{
			desc:      " Create Fail",
			expecterr: errors.DB{Err: errors.Error("new")},
			inputMov: &models.Movie{Name: "Silicon valley",
				Genre: "Comedy", Rating: 4.5,
				ReleaseDate: date, UpdatedAt: dateTime,
				CreatedAt: dateTime, Plot: "Richard, a programmer, " +
					"creates an app called the Pied Piper " +
					"and tries to get investors" +
					" for it. Meanwhile, five other programmers " +
					"struggle to make their " +
					"mark in Silicon Valley.", Released: true},
			mockCall: []interface{}{mock.ExpectExec("insert into movie (name,genre,rating,releaseDate"+
				",updatedAt,createdAt,plot,released) "+
				"VALUES (?, ?, ?,?,?,?,?,?)").WithArgs("Silicon valley",
				"Comedy", 4.5, date, sqlmock.AnyArg(),
				sqlmock.AnyArg(), "Richard, a programmer, creates "+
					"an app called the Pied Piper and tries to get "+
					"investors for it. Meanwhile, five other programmers "+
					"struggle to make their mark in Silicon"+
					" Valley.", true).WillReturnError(errors.Error("new")),
				mock.ExpectQuery("select id , name ,genre" +
					" ,rating , releaseDate ,updatedAt," +
					"createdAt,plot,released from movie" +
					" where deletedAt IS NULL and id=?").
					WithArgs(3).WillReturnError(errors.Error("new")),
			},
		},
	}
	for _, tc := range tests {
		s := New()
		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()
		_, err := s.Create(ctx, tc.inputMov)

		if !reflect.DeepEqual(err, tc.expecterr) {
			t.Errorf("Expected: %v, Got: %v", tc.expecterr, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"ID", "Name",
		"Genre", "Rating", "ReleaseDate",
		"UpdatedAt", "CreatedAt", "Plot",
		"Released"}).AddRow(1, "Silicon valley",
		"Comedy", 4.5, date, dateTime, dateTime,
		"Richard, a programmer, creates an app"+
			" called the Pied Piper and tries to"+
			" get investors for it. Meanwhile, "+
			"five other programmers struggle"+
			" to make their mark in Silicon Valley.", true)

	testCases := []struct {
		mockQuery     interface{}
		expectedError error
		output        []*models.Movie
	}{

		{
			mockQuery: mock.ExpectQuery("SELECT id," +
				"name,genre,rating,releaseDate," +
				"updatedAt,createdAt,plot,released " +
				"FROM movie where deletedAt" +
				" IS NULL").WithArgs().WillReturnRows(rows),
			expectedError: nil,
		},
		{
			mockQuery: mock.ExpectQuery("SELECT id," +
				"name,genre,rating,releaseDate,updatedAt," +
				"createdAt,plot,released FROM " +
				"movie where deletedAt IS NULL").WithArgs().
				WillReturnError(errors.Error("new")),
			expectedError: errors.DB{Err: errors.Error("new")},
		},
	}

	for _, tc := range testCases {
		s := New()
		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()
		_, err := s.GetAll(ctx)

		if !reflect.DeepEqual(err, tc.expectedError) {
			t.Errorf("Expected: %v, Got: %v", tc.expectedError, err)
		}
	}
}
