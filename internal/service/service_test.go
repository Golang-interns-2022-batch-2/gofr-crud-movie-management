package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
	"github.com/iamkakashi/movie-gofr/internal/model"
	"github.com/iamkakashi/movie-gofr/internal/store"
	"gopkg.in/guregu/null.v4"
)

func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	currenttime := time.Now()
	testcase := []struct {
		desc   string
		output []*model.Movie
		erout  error
		mock   *gomock.Call
	}{
		{
			desc: "Success Case",
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
			mock: mockMovieStore.EXPECT().Get(ctx).Return([]*model.Movie{{
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
			}, nil),
		},
		{
			desc:   "Fail Case",
			output: nil,
			erout:  errors.Error("internal server error"),
			mock:   mockMovieStore.EXPECT().Get(ctx).Return(nil, errors.Error("internal server error")),
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			out, err := movieService.Get(ctx)
			if err != nil && !reflect.DeepEqual(err, tc.erout) {
				t.Errorf("Expected %v got %v", tc.erout, err)
			}
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, out)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	currenttime := time.Now()
	testcase := []struct {
		desc   string
		input  int
		output *model.Movie
		erout  error
		mock   *gomock.Call
	}{
		{
			desc:   "Negative id Case",
			input:  -4,
			output: nil,
			erout:  errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc:   "internal server error",
			input:  4,
			output: nil,
			erout:  errors.Error("internal server error"),
			mock:   mockMovieStore.EXPECT().GetByID(ctx, 4).Return(nil, errors.Error("internal server error")),
		},
		{
			desc:  "Success Case",
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
			mock: mockMovieStore.EXPECT().GetByID(ctx, 4).Return(&model.Movie{
				ID:          4,
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			}, nil),
		},
		{
			desc:   "movie does not exist",
			input:  4,
			output: nil,
			erout:  errors.Error("movie not found"),
			mock:   mockMovieStore.EXPECT().GetByID(ctx, 4).Return(nil, errors.Error("movie not found")),
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			out, err := movieService.GetByID(ctx, tc.input)
			if err != nil && !reflect.DeepEqual(err, tc.erout) {
				t.Errorf("Expected %v got %v", tc.erout, err)
			}
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, out)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	currenttime := time.Now()
	testcase := []struct {
		desc   string
		input  *model.Movie
		output *model.Movie
		erout  error
		mock   *gomock.Call
	}{
		{
			desc: "Success Case",
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
			mock: mockMovieStore.EXPECT().Create(ctx, &model.Movie{
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
			}).Return(&model.Movie{
				ID:          4,
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			}, nil),
		},
		{
			desc: "Invalid Name",
			input: &model.Movie{
				Name:        null.NewString("", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: nil,
			erout:  errors.InvalidParam{},
		},
		{
			desc: "Fail Case",
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
			erout: errors.Error("failed to create movie"),
			mock: mockMovieStore.EXPECT().Create(ctx, &model.Movie{
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			}).Return(nil, errors.Error("failed to create movie")),
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			out, err := movieService.Create(ctx, tc.input)
			if err != nil && !reflect.DeepEqual(err, tc.erout) {
				t.Errorf("Expected %v got %v", tc.erout, err)
			}
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, out)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	currenttime := time.Now()
	testcase := []struct {
		desc   string
		input  *model.Movie
		output *model.Movie
		erout  error
		mock   *gomock.Call
	}{
		{
			desc: "Success Case",
			input: &model.Movie{
				Name:        null.NewString("Deadpool2", true),
				Genre:       null.NewString("Comedy", true),
				Rating:      null.NewFloat(4.0, true),
				ReleaseDate: null.NewString("2018-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: &model.Movie{
				ID:          4,
				Name:        null.NewString("Deadpool2", true),
				Genre:       null.NewString("Comedy", true),
				Rating:      null.NewFloat(4.0, true),
				ReleaseDate: null.NewString("2018-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			erout: nil,
			mock: mockMovieStore.EXPECT().Update(ctx, &model.Movie{
				Name:        null.NewString("Deadpool2", true),
				Genre:       null.NewString("Comedy", true),
				Rating:      null.NewFloat(4.0, true),
				ReleaseDate: null.NewString("2018-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("superhero movie", true),
				Released:    null.NewBool(true, true),
			}).Return(&model.Movie{
				ID:          4,
				Name:        null.NewString("Deadpool2", true),
				Genre:       null.NewString("Comedy", true),
				Rating:      null.NewFloat(4.0, true),
				ReleaseDate: null.NewString("2018-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("superhero movie", true),
				Released:    null.NewBool(true, true),
			}, nil),
		},
		{
			desc: "Invalid Rating",
			input: &model.Movie{
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(6.0, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: nil,
			erout:  errors.InvalidParam{},
		},
		{
			desc: "Invalid Name",
			input: &model.Movie{
				Name:        null.NewString("", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			},
			output: nil,
			erout:  errors.InvalidParam{},
		},
		{
			desc: "Fail Case",
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
			mock: mockMovieStore.EXPECT().Update(ctx, &model.Movie{
				Name:        null.NewString("Deadpool", true),
				Genre:       null.NewString("Comedy, Action", true),
				Rating:      null.NewFloat(4.8, true),
				ReleaseDate: null.NewString("2014-12-17", true),
				UpdatedAt:   null.NewTime(currenttime, true),
				CreatedAt:   null.NewTime(currenttime, true),
				Plot:        null.NewString("This is a superhero movie", true),
				Released:    null.NewBool(true, true),
			}).Return(nil, errors.Error("failed to create movie")),
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			out, err := movieService.Update(ctx, tc.input)
			if err != nil && !reflect.DeepEqual(err, tc.erout) {
				t.Errorf("Expected %v got %v", tc.erout, err)
			}
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, out)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := gofr.NewContext(nil, nil, nil)
	ctx.Context = context.Background()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	testcase := []struct {
		desc  string
		input int
		erout error
		mock  *gomock.Call
	}{
		{
			desc:  "Negative id Case",
			input: -4,
			erout: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc:  "Success Case",
			input: 4,
			erout: nil,
			mock:  mockMovieStore.EXPECT().Delete(ctx, 4).Return(nil),
		},
		{
			desc:  "Fail Case",
			input: 4,
			erout: errors.Error("failed to delete"),
			mock:  mockMovieStore.EXPECT().Delete(ctx, 4).Return(errors.Error("failed to delete")),
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			err := movieService.Delete(ctx, tc.input)
			if err != nil && !reflect.DeepEqual(err, tc.erout) {
				t.Errorf("Expected %v got %v", tc.erout, err)
			}
		})
	}
}
