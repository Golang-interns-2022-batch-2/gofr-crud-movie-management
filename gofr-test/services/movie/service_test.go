package movie

import (
	"context"
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"

	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/models"

	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/datastore"

	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := datastore.NewMockMovie(ctrl)
	s := New(mockService)

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testcases := []struct {
		desc     string
		id       int
		expected *models.Movie
		mockCall *gomock.Call
	}{
		{desc: "Get success",
			id: 3,
			expected: &models.Movie{ID: 3,
				Name:  "Silicon valley",
				Genre: "Comedy", Rating: 4.5,
				ReleaseDate: date, UpdatedAt: dateTime,
				CreatedAt: dateTime, Plot: "Richard, a programmer, " +
					"creates an app called the" +
					" Pied Piper and tries to get " +
					"investors for it. Meanwhile," +
					" five other programmers struggle" +
					" to make their mark in " +
					"Silicon Valley.", Released: true},
			mockCall: mockService.EXPECT().GetByID(ctx, 3).Return(&models.Movie{ID: 3,
				Name:        "Silicon valley",
				Genre:       "Comedy",
				Rating:      4.5,
				ReleaseDate: date,
				UpdatedAt:   dateTime,
				CreatedAt:   dateTime,
				Plot: "Richard, a programmer, " +
					"creates an app called the Pied " +
					"Piper and tries to get investors" +
					" for it. Meanwhile, five other " +
					"programmers struggle to make their" +
					" mark in Silicon Valley.", Released: true}, nil),
		},
		{
			desc:     "Get Fail",
			id:       -1,
			expected: nil,
			mockCall: nil,
		},
	}

	for i, tc := range testcases {
		res, _ := s.GetByID(ctx, tc.id)
		if !reflect.DeepEqual(res, tc.expected) {
			t.Errorf("%v [TEST%d]Failed. Got %v\tExpected %v\n", tc.desc, i+1, res, tc.expected)
		}
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := datastore.NewMockMovie(ctrl)
	s := New(mockService)

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()

	tests := []struct {
		desc     string
		id       int
		err      error
		mockCall *gomock.Call
	}{
		{desc: "Delete success", id: 3, err: nil, mockCall: mockService.EXPECT().Delete(ctx, 3).Return(nil)},
		{desc: "Delete fail",
			id:  -1,
			err: errors.InvalidParam{Param: []string{"id"}},
			//mockCall: mockService.EXPECT().Delete(ctx, -1).Return(errors.New("negative id"))},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := s.Delete(ctx, tc.id)

			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("Expected: %v, Got: %v", tc.err, err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := datastore.NewMockMovie(ctrl)
	s := New(mockService)

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testMovie := models.Movie{ID: 3,
		Name:        "Silicon valley",
		Genre:       "Comedy",
		Rating:      4.5,
		ReleaseDate: date,
		UpdatedAt:   dateTime,
		CreatedAt:   dateTime,
		Plot: "Richard, a programmer, " +
			"creates an app called " +
			"the Pied Piper and tries " +
			"to get investors for it. " +
			"Meanwhile, five other programmers" +
			" struggle to make their " +
			"mark in Silicon Valley.", Released: true}

	testMovie1 := models.Movie{ID: -1,
		Name:        "Silicon valley",
		Genre:       "Comedy",
		Rating:      4.5,
		ReleaseDate: date,
		UpdatedAt:   dateTime,
		CreatedAt:   dateTime,
		Plot: "Richard, a programmer, " +
			"creates an app called " +
			"the Pied Piper and tries " +
			"to get investors for it. " +
			"Meanwhile, five other programmers" +
			" struggle to make their " +
			"mark in Silicon Valley.", Released: true}

	tests := []struct {
		desc          string
		id            int
		mv            models.Movie
		expectedError error
		mockCall      []*gomock.Call
	}{
		{
			desc: "Success",
			id:   3,
			mv:   testMovie,
			mockCall: []*gomock.Call{
				mockService.EXPECT().Update(ctx, &testMovie).Return(&testMovie, nil),
			},
		},
		{
			desc:          "Failure",
			id:            -1,
			mv:            testMovie1,
			expectedError: errors.InvalidParam{Param: []string{"id"}},
			//mockCall: []*gomock.Call{
			//	mockService.EXPECT().Update(ctx, gomock.Any()).Return(nil, sql.ErrNoRows),
			//},
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			_, err := s.Update(ctx, &test.mv)
			if !reflect.DeepEqual(err, test.expectedError) {
				t.Errorf("expected: %v, Got: %v", test.expectedError, err)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := datastore.NewMockMovie(ctrl)
	s := New(mockService)

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testMovie := models.Movie{ID: 3,
		Name:        "Silicon valley",
		Genre:       "Comedy",
		Rating:      4.5,
		ReleaseDate: date,
		UpdatedAt:   dateTime,
		CreatedAt:   dateTime,
		Plot: "Richard, a programmer, " +
			"creates an app called the Pied Piper" +
			" and tries to get investors for it. Meanwhile, " +
			"five other programmers struggle to " +
			"make their mark in Silicon Valley.",
		Released: true}

	testMovie2 := models.Movie{ID: 3,
		Name:        "Silicon valley",
		Genre:       "",
		Rating:      4.5,
		ReleaseDate: date,
		UpdatedAt:   dateTime,
		CreatedAt:   dateTime,
		Plot: "Richard, a programmer, " +
			"creates an app called the " +
			"Pied Piper and tries to get " +
			"investors for it. Meanwhile," +
			" five other programmers struggle" +
			" to make their mark in" +
			" Silicon Valley.", Released: true}

	testMovie1 := models.Movie{ID: 11,
		Name:        "",
		Genre:       "Comedy",
		Rating:      0.0,
		ReleaseDate: date,
		UpdatedAt:   dateTime,
		CreatedAt:   dateTime,
		Plot: "Richard, a programmer," +
			" creates an app called the Pied " +
			"Piper and tries to get investors " +
			"for it. Meanwhile, five other " +
			"programmers struggle to make " +
			"their mark in Silicon Valley.", Released: true}

	tests := []struct {
		desc     string
		output   models.Movie
		error    error
		mockcall *gomock.Call
	}{
		{
			desc: "Create success", output: testMovie, error: nil,
			mockcall: mockService.EXPECT().Create(ctx, &testMovie).Return(&testMovie, nil),
		},
		{
			desc:   "Create fail",
			output: models.Movie{},
			error:  errors.InvalidParam{},
		},
		{
			desc:   "Create fail",
			output: testMovie1,
			error:  errors.InvalidParam{},
		},
		{
			desc:   "Create fail3",
			output: testMovie2,
			error:  errors.InvalidParam{},
		},
	}

	for i, tc := range tests {
		_, err := s.Create(ctx, &tc.output)

		if !reflect.DeepEqual(err, tc.error) {
			t.Errorf(" [TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, tc.error)
		}
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := datastore.NewMockMovie(ctrl)
	s := New(mockService)

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()

	tests := []struct {
		id       int
		error    error
		mockcall []*gomock.Call
	}{
		{
			id: 3, error: nil,
			mockcall: []*gomock.Call{
				mockService.EXPECT().GetAll(ctx).Return([]*models.Movie{}, nil),
			},
		},
		{
			id: -1, error: sql.ErrNoRows,
			mockcall: []*gomock.Call{
				mockService.EXPECT().GetAll(ctx).Return(nil, sql.ErrNoRows),
			},
		},
	}

	for i, tc := range tests {
		_, err := s.GetAll(ctx)

		if !reflect.DeepEqual(err, tc.error) {
			t.Errorf(" [TEST%d]Failed. Got %v\tExpected %v\n", i+1, err, tc.error)
		}
	}
}
