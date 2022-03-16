package moviestore

import (
	"context"
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/datastore"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/DATA-DOG/go-sqlmock"
	"golangprog/gofr-test/Internal/models"
	"reflect"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	date, _ := time.Parse(time.RFC3339, "2014-12-17")
	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testCases := []struct {
		id          int
		mockQuery   interface{}
		expectError error
		desc        string
	}{
		{
			id: 1,
			mockQuery: mock.ExpectQuery("Select id, name, genre, rating, releaseDate, updatedAt, " +
				"createdAt, plot, released FROM movie WHERE id = ? and deletedAt " +
				"IS NULL").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"Id", "Name",
				"Genre", "Rating", "ReleaseDate", "UpdateAt", "CreatedAt", "Plot", "Released"}).AddRow(1, "Silicon Valley",
				"Comedy", 4.5, date, dateTime, dateTime, `"Richard, a programmer, creates an app called the Pied Piper 
			and tries to get investors for it. Meanwhile, five other programmers struggle to make their mark in Silicon Valley."`, true)),
			expectError: nil,
			desc:        "success case",
		},
		{
			id: -1,
			mockQuery: mock.ExpectQuery("Select id, name, genre, rating, releaseDate, updatedAt, " +
				"createdAt, plot, released FROM movie WHERE id = ? and deletedAt IS NULL").WithArgs(-1).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
			desc:        "failure case",
		},
	}

	for _, testCase := range testCases {
		s := New()
		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()
		_, err := s.GetByID(ctx, testCase.id)

		if !reflect.DeepEqual(err, testCase.expectError) {
			t.Errorf("expected: %v, got: %v", testCase.expectError, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	date, _ := time.Parse(time.RFC3339, "2014-12-17")
	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testCases := []struct {
		id          int
		mockQuery   interface{}
		expectError error
		desc        string
	}{
		{
			id: 1,
			mockQuery: mock.ExpectQuery("Select id, name, genre, rating, releaseDate, updatedAt, createdAt, plot, " +
				"released FROM movie WHERE deletedAt IS NULL").WithArgs().WillReturnRows(sqlmock.NewRows([]string{"Id", "Name", "Genre", "Rating",
				"ReleaseDate", "UpdateAt", "CreatedAt", "Plot", "Released"}).AddRow(6, "Silicon Valley", "Comedy", 4.5, date, dateTime,
				dateTime, "Richard, a programmer, creates an app called the Pied Piper and tries to get investors for it. "+
					"Meanwhile, five other programmers struggle to make their mark in Silicon Valley.", true)),
			expectError: nil,
			desc:        "success case",
		},
		{
			id: 3,
			mockQuery: mock.ExpectQuery("Select id, name, genre, rating, releaseDate, updatedAt, createdAt, plot, released FROM movie " +
				"WHERE deletedAt IS NULL").WithArgs().WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
			desc:        "failure case",
		},
		{
			id: -1,
			mockQuery: mock.ExpectQuery("Select id, name, genre, rating, releaseDate, updatedAt, createdAt, plot, released FROM movie " +
				"WHERE deletedAt IS NULL").WithArgs().WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
			desc:        "failure case",
		},
	}

	for _, testCase := range testCases {
		s := New()
		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()
		return
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := s.GetAll(ctx)

			if !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected: %v, got: %v", testCase.expectError, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testCases := []struct {
		id          int
		mockQuery   *sqlmock.ExpectedExec
		expectError error
		desc        string
	}{
		{
			id:          1,
			mockQuery:   mock.ExpectExec("update movie set deletedAt = Now() Where id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1)),
			expectError: nil,
			desc:        "success case",
		},
		{
			id:          3,
			mockQuery:   mock.ExpectExec("update movie set deletedAt = Now() Where id = ?").WithArgs(3).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
			desc:        "failure case",
		},
		{
			id:          -1,
			mockQuery:   mock.ExpectExec("update movie set deletedAt = Now() Where id = ?").WithArgs(-1).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
			desc:        "failure case",
		},
	}
	for _, testCase := range testCases {
		s := New()
		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()
		return
		t.Run(testCase.desc, func(t *testing.T) {
			err := s.Delete(ctx, testCase.id)

			if !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected: %v, got: %v", testCase.expectError, err)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	date, _ := time.Parse(time.RFC3339, "2014-12-17")
	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testCases := []struct {
		id           int
		mockQuery    *sqlmock.ExpectedExec
		expectOutput *models.Movie
		expectError  error
		desc         string
	}{
		{
			id: 1,
			mockQuery: mock.ExpectExec("INSERT into movie (name, genre, rating, "+
				"releaseDate, createdAt, updatedAt, plot, released) VALUES (?,?,?,?,?,?,?,?)").WithArgs("Silicon Valley", "Comedy", 4.5, date,
				sqlmock.AnyArg(), sqlmock.AnyArg(), "Richard, a programmer, creates an app called the Pied Piper and tries to get investors "+
					"for it. Meanwhile, five other programmers struggle "+
					"to make their mark in Silicon Valley.", true).WillReturnResult(sqlmock.NewResult(1, 1)),
			expectOutput: &models.Movie{ID: 1, Name: "Silicon Valley", Genre: "Comedy", Rating: 4.5, ReleaseDate: date, CreatedAt: dateTime,
				UpdatedAt: dateTime, Plot: "Richard, a programmer, creates an app called the Pied Piper " +
					"and tries to get investors for it. Meanwhile, five other programmers struggle" +
					" to make their mark in Silicon Valley.", Released: true},
			expectError: nil,
			desc:        "success case",
		},
		{
			id: 3,
			mockQuery: mock.ExpectExec("INSERT into movie (name, genre, rating, "+
				"releaseDate, createdAt, updatedAt, plot, released) VALUES (?,?,?,?,?,?,?,?)").WithArgs("", "", 0.0, time.Time{}, sqlmock.AnyArg(),
				sqlmock.AnyArg(), "", false).WillReturnError(sql.ErrNoRows),
			expectOutput: &models.Movie{ID: 3, Name: "", Genre: "", Rating: 0.0, ReleaseDate: time.Time{},
				CreatedAt: time.Time{}, UpdatedAt: time.Time{}, Plot: "", Released: false},
			expectError: sql.ErrNoRows,
			desc:        "failure case",
		},
		{
			id: -1,
			mockQuery: mock.ExpectExec("INSERT into movie (name, genre, rating, "+
				"releaseDate, createdAt, updatedAt, plot, released) VALUES (?,?,?,?,?,?,?,?)").WithArgs("", "", 0.0, time.Time{},
				sqlmock.AnyArg(), sqlmock.AnyArg(), "", false).WillReturnError(sql.ErrNoRows),
			expectError: sql.ErrNoRows,
			expectOutput: &models.Movie{ID: -1, Name: "", Genre: "", Rating: 0.0, ReleaseDate: time.Time{}, CreatedAt: time.Time{},
				UpdatedAt: time.Time{}, Plot: "", Released: false},
			desc: "failure case",
		},
	}

	for _, testCase := range testCases {
		s := New()
		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()
		return
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := s.Create(ctx, testCase.expectOutput)

			if !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected: %v, got: %v", testCase.expectError, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	date, _ := time.Parse(time.RFC3339, "2014-12-17")
	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testCases := []struct {
		id          int
		mockQuery   *sqlmock.ExpectedExec
		inp         models.Movie
		expectError error
		desc        string
	}{
		{
			id: 1,
			mockQuery: mock.ExpectExec("Update movie SET name=?, genre=?, rating=?, releaseDate=?, updatedAt=?, createdAt=?, plot=?, "+
				"released=? WHERE id = ? and deletedAt IS NULL").WithArgs("Silicon Valley", "Comedy", 4.5, date,
				sqlmock.AnyArg(), sqlmock.AnyArg(), "Richard, a programmer, creates an app called the Pied Piper and "+
					"tries to get investors for it. Meanwhile, five other programmers struggle to make their mark in Silicon "+
					"Valley.", true, 1).WillReturnResult(sqlmock.NewResult(1, 1)),
			inp: models.Movie{ID: 1, Name: "Silicon Valley", Genre: "Comedy", Rating: 4.5, ReleaseDate: date,
				CreatedAt: dateTime, UpdatedAt: dateTime,
				Plot: "Richard, a programmer, creates an app called " +
					"the Pied Piper and tries to get investors for it. Meanwhile, five other programmers " +
					"struggle to make their mark in Silicon Valley.", Released: true},
			expectError: nil,
			desc:        "success case",
		},
		{
			id: 3,
			mockQuery: mock.ExpectExec("Update movie SET name=?, genre=?, rating=?, releaseDate=?, updatedAt=?, createdAt=?, plot=?, released=? "+
				"WHERE id = ? and deletedAt IS NULL").WithArgs("", "", 0.0, time.Time{}, sqlmock.AnyArg(),
				sqlmock.AnyArg(), "", false, 3).WillReturnError(sql.ErrNoRows),
			inp: models.Movie{ID: 3, Name: "", Genre: "", Rating: 0.0, ReleaseDate: time.Time{}, CreatedAt: time.Time{}, UpdatedAt: time.Time{},
				Plot: "", Released: false},
			expectError: sql.ErrNoRows,
			desc:        "failure case",
		},
		{
			id: -1,
			mockQuery: mock.ExpectExec("Update movie SET name=?, genre=?, rating=?, releaseDate=?, updatedAt=?, "+
				"createdAt=?, plot=?, released=? WHERE id = ? and deletedAt IS NULL").WithArgs("", "", 0.0, time.Time{},
				sqlmock.AnyArg(), sqlmock.AnyArg(), "", false, -1).WillReturnError(sql.ErrNoRows),
			inp: models.Movie{ID: -1, Name: "", Genre: "", Rating: 0.0, ReleaseDate: time.Time{}, CreatedAt: time.Time{}, UpdatedAt: time.Time{},
				Plot: "", Released: false},
			expectError: sql.ErrNoRows,
			desc:        "failure case",
		},
	}
	for _, testCase := range testCases {
		s := New()
		ctx := gofr.NewContext(nil, nil, &gofr.Gofr{DataStore: datastore.DataStore{ORM: db}})
		ctx.Context = context.Background()
		return
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := s.Update(ctx, &testCase.inp)

			if !reflect.DeepEqual(err, testCase.expectError) {
				t.Errorf("expected: %v, got: %v", testCase.expectError, err)
			}
		})
	}
}
