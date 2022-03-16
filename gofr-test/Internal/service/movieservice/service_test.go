package movieservice

import (
	"context"
	"database/sql"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"errors"
	"github.com/golang/mock/gomock"
	"golangprog/gofr-test/Internal/models"
	"golangprog/gofr-test/Internal/store"
	"reflect"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := store.NewMockMovie(ctrl)
	s := New(mockService)
	date, _ := time.Parse(time.RFC3339, "2014-12-17")
	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	testMovie := models.Movie{
		ID:          1,
		Name:        "Silicon Valley",
		Genre:       "Comedy",
		Rating:      4.5,
		ReleaseDate: date,
		CreatedAt:   dateTime,
		UpdatedAt:   dateTime,
		Plot: `"Richard, a programmer, creates an app called the Pied Piper and tries to get investors for it. 
		Meanwhile, five other programmers struggle to make their mark in Silicon Valley."`,
		Released: true,
	}

	testCases := []struct {
		desc          string
		id            int
		expectedError error
		mock          []*gomock.Call
	}{
		{desc: "Get success",
			id:            1,
			expectedError: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().GetByID(ctx, 1).Return(&testMovie, nil),
			},
		},
		{desc: "Get failure",
			id:            0,
			expectedError: errors.New("error invalid id"),
			mock:          nil,
		},
		{desc: "Get failure",
			id:            -1,
			expectedError: errors.New("error invalid id"),
			mock:          nil,
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.desc, func(t *testing.T) {

			_, err := s.GetByID(ctx, testCase.id)
			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected: %v, got: %v", testCase.expectedError, err)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := store.NewMockMovie(ctrl)
	s := New(mockService)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	testCases := []struct {
		desc          string
		id            int
		expectedError error
		mock          []*gomock.Call
	}{
		{desc: "Get success",
			id:            1,
			expectedError: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().GetAll(ctx).Return([]*models.Movie{}, nil),
			},
		},
		{desc: "Get failure",
			id:            0,
			expectedError: sql.ErrNoRows,
			mock: []*gomock.Call{
				mockService.EXPECT().GetAll(ctx).Return(nil, sql.ErrNoRows),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := s.GetAll(ctx)
			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected: %v, got: %v", testCase.expectedError, err)
			}
		})
	}
}

func TestMovieDel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := store.NewMockMovie(ctrl)
	s := New(mockService)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	testCases := []struct {
		desc          string
		id            int
		expectedError error
		mock          []*gomock.Call
	}{
		{desc: "Get success",
			id:            1,
			expectedError: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().Delete(ctx, 1).Return(nil),
			},
		},
		{desc: "Get failure",
			id:            0,
			expectedError: errors.New("error invalid id"),
			mock: []*gomock.Call{
				mockService.EXPECT().Delete(ctx, 0).Return(errors.New("error invalid id")),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			err := s.Delete(ctx, testCase.id)
			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected: %v, got: %v", testCase.expectedError, err)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := store.NewMockMovie(ctrl)
	s := New(mockService)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	date, _ := time.Parse(time.RFC3339, "2014-12-17")
	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testMovie := models.Movie{
		ID:          1,
		Name:        "Silicon Valley",
		Genre:       "Comedy",
		Rating:      4.5,
		ReleaseDate: date,
		CreatedAt:   dateTime,
		UpdatedAt:   dateTime,
		Plot: "Richard, a programmer, creates an app called the Pied Piper and tries to get investors for it. " +
			"Meanwhile, five other programmers struggle to make their mark in Silicon Valley.",
		Released: true,
	}
	testMovie1 := models.Movie{
		ID:          0,
		Name:        "",
		Genre:       "",
		Rating:      4.5,
		ReleaseDate: date,
		CreatedAt:   dateTime,
		UpdatedAt:   dateTime,
		Plot: "Richard, a programmer, creates an app called the Pied Piper and tries to get investors for it. " +
			"Meanwhile, five other programmers struggle to make their mark in Silicon Valley.",
		Released: true,
	}
	testCases := []struct {
		desc          string
		id            int
		testMv        models.Movie
		expectedError error
		mock          []*gomock.Call
	}{
		{desc: "Get success",
			id:            1,
			testMv:        testMovie,
			expectedError: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().Create(ctx, &testMovie).Return(&testMovie, nil),
			},
		},
		{desc: "Get failure",
			id:            0,
			testMv:        testMovie1,
			expectedError: errors.New("error invalid name"),
			mock:          nil,
		},
	}

	for _, testCase := range testCases {
		_, err := s.Create(ctx, &testCase.testMv)
		if !reflect.DeepEqual(err, testCase.expectedError) {
			t.Errorf("expected: %v, got: %v", testCase.expectedError, err)
		}
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := store.NewMockMovie(ctrl)
	s := New(mockService)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	date, _ := time.Parse(time.RFC3339, "2014-12-17")
	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testMovie := models.Movie{
		ID:          1,
		Name:        "Silicon Valley",
		Genre:       "Comedy",
		Rating:      4.5,
		ReleaseDate: date,
		CreatedAt:   dateTime,
		UpdatedAt:   dateTime,
		Plot: "Richard, a programmer, creates an app called the Pied Piper and tries to get investors for it. " +
			"Meanwhile, five other programmers struggle to make their mark in Silicon Valley.",
		Released: true,
	}
	testMovie1 := models.Movie{
		ID:          0,
		Name:        "",
		Genre:       "",
		Rating:      4.5,
		ReleaseDate: date,
		CreatedAt:   dateTime,
		UpdatedAt:   dateTime,
		Plot: "Richard, a programmer, creates an app called the Pied Piper and tries to get investors for it. " +
			"Meanwhile, five other programmers struggle to make their mark in Silicon Valley.",
		Released: true,
	}

	testCases := []struct {
		desc          string
		id            int
		mv            models.Movie
		expectedError error
		mock          []*gomock.Call
	}{
		{desc: "Get success",
			id:            1,
			expectedError: nil,
			mv:            testMovie,
			mock: []*gomock.Call{
				mockService.EXPECT().Update(ctx, &testMovie).Return(&testMovie, nil),
			},
		},
		{desc: "Get failure",
			id:            0,
			expectedError: errors.New("error invalid name"),
			mv:            testMovie1,
			mock:          nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			_, err := s.Update(ctx, &testCase.mv)
			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Errorf("expected: %v, got: %v", testCase.expectedError, err)
			}
		})
	}
}
