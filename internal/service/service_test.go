package service

import (
	"reflect"
	"testing"
	"time"

	gofrerr "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/golang/mock/gomock"
	"github.com/shivam/Crud_Gofr/internal/models"
	"github.com/shivam/Crud_Gofr/internal/store"
)

func Test_GetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)

	CurrTime := time.Now()
	testcase := []struct {
		desc   string
		output []*models.Movie
		err    error
		query  *gomock.Call
	}{
		{
			desc: "Success Case",
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
			query: mockMovieStore.EXPECT().GetAll(gomock.Any()).Return([]*models.Movie{
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
			}, nil),
		},
		{
			desc:   "Fail Case",
			output: nil,
			err:    gofrerr.Error("error1 in service layer while calling GetAll function"),
			query:  mockMovieStore.EXPECT().GetAll(gomock.Any()).Return(nil, gofrerr.Error("error1 in service layer while calling GetAll function")),
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		ctx := gofr.NewContext(nil, nil, gofr.New())

		t.Run("", func(t *testing.T) {
			out, err := movieService.GetAllService(ctx)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected %v got %v", tc.err, err)
			}
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, out)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	ctx := gofr.NewContext(nil, nil, gofr.New())

	CurrTime := time.Now()
	testcase := []struct {
		desc   string
		input  int
		output *models.Movie
		err    error
		query  *gomock.Call
	}{
		{
			desc:   "Fail case - Negative Id ",
			input:  -4,
			output: nil,
			err:    gofrerr.Error("id cannot be nagative or zero"),
		},
		{
			desc:   "Testing call function in service layer",
			input:  4,
			output: nil,
			err:    gofrerr.Error("error in service layer while calling getby id"),
			query:  mockMovieStore.EXPECT().GetByID(ctx, 4).Return(nil, gofrerr.Error("error in service layer while calling getby id")),
		},
		{
			desc:  "Success Case",
			input: 4,
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
			query: mockMovieStore.EXPECT().GetByID(ctx, 4).Return(&models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			}, nil),
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			out, err := movieService.GetByIDService(ctx, tc.input)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected %v got %v", tc.err, err)
			}
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, out)
			}
		})
	}
}

func Test_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	ctx := gofr.NewContext(nil, nil, gofr.New())

	CurrTime := time.Now()
	testcase := []struct {
		desc   string
		input  *models.Movie
		output *models.Movie
		err    error
		query  *gomock.Call
	}{
		{
			desc: "Success Case",
			input: &models.Movie{
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
			query: mockMovieStore.EXPECT().Create(ctx, &models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			}).Return(&models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			}, nil),
		},
		{
			desc: "Invalid Name",
			input: &models.Movie{
				ID:          1,
				Name:        "",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			},
			output: nil,
			err:    gofrerr.Error("name cannot be null"),
		},
		{
			desc: "Invalid Genre",
			input: &models.Movie{
				ID:          1,
				Name:        "Name",
				Genre:       "",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			},
			output: nil,
			err:    gofrerr.Error("genre cannot be null"),
		},
		{
			desc: "Invalid Plot",
			input: &models.Movie{
				ID:          1,
				Name:        "Name",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "",
				Released:    true,
			},
			output: nil,
			err:    gofrerr.Error("plot cannot be empty"),
		},
		{
			desc: "Invalid Rating",
			input: &models.Movie{
				ID:          1,
				Name:        "Name",
				Genre:       "Action",
				Rating:      -2,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "plot",
				Released:    true,
			},
			output: nil,
			err:    gofrerr.Error("rating cannot be negative"),
		},

		{
			desc: "Invalid ID",
			input: &models.Movie{
				ID:          -11,
				Name:        "Name",
				Genre:       "Action",
				Rating:      2.5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "plot",
				Released:    true,
			},
			output: nil,
			err:    gofrerr.Error("id cannot be negative"),
		},

		{
			desc: "Fail Case",
			input: &models.Movie{
				ID:          1,
				Name:        "Name",
				Genre:       "Action",
				Rating:      2.5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "plot",
				Released:    true,
			},
			output: nil,
			err:    gofrerr.Error("error in service layer while caaling create function"),
			query: mockMovieStore.EXPECT().Create(ctx, &models.Movie{
				ID:          1,
				Name:        "Name",
				Genre:       "Action",
				Rating:      2.5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "plot",
				Released:    true,
			}).Return(nil, gofrerr.Error("error in service layer while caaling create function")),
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			out, err := movieService.InsertService(ctx, tc.input)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected %v got %v", tc.err, err)
			}
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, out)
			}
		})
	}
}

func Test_Delete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	ctx := gofr.NewContext(nil, nil, gofr.New())
	CurrTime := time.Now()
	testcase := []struct {
		desc  string
		input int
		err   error
		query []interface{}
	}{
		{
			desc:  "Success Case",
			input: 2,
			err:   nil,
			query: []interface{}{mockMovieStore.EXPECT().GetByID(ctx, 2).Return(&models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			}, nil), mockMovieStore.EXPECT().DeleteByID(ctx, 2).Return(nil)},
		},
		{
			desc:  "Fail - Negative id ",
			input: -1,
			err:   gofrerr.Error("id cannot be nagative"),
		},
		{
			desc:  "Fail Case",
			input: 2,
			err:   gofrerr.Error("error in service layer for deleting id"),
			query: []interface{}{mockMovieStore.EXPECT().GetByID(ctx, 2).Return(&models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			}, nil), mockMovieStore.EXPECT().DeleteByID(ctx, 2).Return(gofrerr.Error("error in service layer for deleting id"))},
		},

		{
			desc:  "Failed to find id in db ",
			input: 1,
			err:   gofrerr.Error("id no present in the database"),
			query: []interface{}{mockMovieStore.EXPECT().GetByID(ctx, 1).Return(nil, gofrerr.Error("id no present in the database"))},
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			err := movieService.DeleteService(ctx, tc.input)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected %v got %v", tc.err, err)
			}
		})
	}
}
func Test_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockMovieStore := store.NewMockMovieStorer(mockCtrl)
	ctx := gofr.NewContext(nil, nil, gofr.New())

	CurrTime := time.Now()
	testcase := []struct {
		desc   string
		input  *models.Movie
		output *models.Movie
		err    error
		query  interface{}
	}{
		{
			desc: "Success Case",
			input: &models.Movie{
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
			query: []interface{}{mockMovieStore.EXPECT().GetByID(ctx, 1).Return(&models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			}, nil),
				mockMovieStore.EXPECT().Update(ctx, 1, &models.Movie{
					ID:          1,
					Name:        "MazeRunner",
					Genre:       "Action",
					Rating:      5,
					ReleaseDate: "2022-12-17",
					UpdatedAt:   CurrTime.String(),
					CreatedAt:   CurrTime.String(),
					Plot:        "PLOT!",
					Released:    true,
				}).Return(&models.Movie{
					ID:          1,
					Name:        "MazeRunner",
					Genre:       "Action",
					Rating:      5,
					ReleaseDate: "2022-12-17",
					UpdatedAt:   CurrTime.String(),
					CreatedAt:   CurrTime.String(),
					Plot:        "PLOT!",
					Released:    true,
				}, nil)},
		},

		{
			desc: "Invalid ID",
			input: &models.Movie{
				ID:          -1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			},
			output: nil,
			err:    gofrerr.Error("ID cannot be negative"),
		},

		{
			desc: "Id not in db",
			input: &models.Movie{
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
			output: nil,
			err:    gofrerr.Error("id no present in the database"),
			query: []interface{}{mockMovieStore.EXPECT().GetByID(ctx, 1).
				Return(nil, gofrerr.Error("error in service layer while calling getby id"))},
		},

		{
			desc: "Id not in db",
			input: &models.Movie{
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
			output: nil,
			err:    gofrerr.Error("error while calling update in service layer"),
			query: []interface{}{mockMovieStore.EXPECT().GetByID(ctx, 1).Return(&models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			}, nil), mockMovieStore.EXPECT().Update(ctx, 1, &models.Movie{
				ID:          1,
				Name:        "MazeRunner",
				Genre:       "Action",
				Rating:      5,
				ReleaseDate: "2022-12-17",
				UpdatedAt:   CurrTime.String(),
				CreatedAt:   CurrTime.String(),
				Plot:        "PLOT!",
				Released:    true,
			}).Return(nil, gofrerr.Error("error while calling update in service layer"))},
		},
	}

	movieService := New(mockMovieStore)

	for _, tc := range testcase {
		tc := tc

		t.Run("", func(t *testing.T) {
			out, err := movieService.UpdatedService(ctx, tc.input)
			if err != nil && err.Error() != tc.err.Error() {
				t.Errorf("Expected %v got %v", tc.err, err)
			}
			if !reflect.DeepEqual(out, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, out)
			}
		})
	}
}
