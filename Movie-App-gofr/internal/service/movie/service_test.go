package movie

import (
	"database/sql"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/model"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/store"
	"github.com/golang/mock/gomock"
)

func TestService_CreateMovie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieStore := store.NewMockMovieInterface(ctrl)
	tcs := []struct {
		mock        *gomock.Call
		expectedErr error
		movie       model.MovieModel
	}{
		{
			movie: model.MovieModel{
				ID:          1,
				Name:        "",
				Genre:       "comedy",
				Rating:      4.5,
				ReleaseDate: "2014-12-24",
				UpdatedAt:   "2022-03-06 19:29:33",
				CreatedAt:   "2022-03-06 13:53:00",
				DeletedAt:   "2022-03-06 13:53:00",
				Plot:        "2022-03-06 13:53:00",
				Released:    true,
			},

			expectedErr: &errors.Response{
				StatusCode: 400,
				Code:       "INVLD_NAME",
				Reason:     "name is empty",
			},
		},
		{
			movie: model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "",
				Rating:      4.5,
				ReleaseDate: "2014-12-24",
				UpdatedAt:   "2022-03-06 19:29:33",
				CreatedAt:   "2022-03-06 13:53:00",
				DeletedAt:   "2022-03-06 13:53:00",
				Plot:        "2022-03-06 13:53:00",
				Released:    true,
			},
			expectedErr: &errors.Response{
				StatusCode: 400,
				Code:       "INVLD_GENRE",
				Reason:     "genre is empty",
			},
		},
		{
			movie: model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "comedy",
				Rating:      4.5,
				ReleaseDate: "2014-12-24",
				UpdatedAt:   "2022-03-06 19:29:33",
				CreatedAt:   "2022-03-06 13:53:00",
				DeletedAt:   "2022-03-06 13:53:00",
				Plot:        "2022-03-06 13:53:00",
				Released:    true,
			},
			mock:        movieStore.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(nil, errors.Error("Server Error")),
			expectedErr: errors.Error("Server Error"),
		},
		{
			movie: model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "comedy",
				Rating:      4.5,
				ReleaseDate: "2014-12-24",
				UpdatedAt:   "2022-03-06 19:29:33",
				CreatedAt:   "2022-03-06 13:53:00",
				DeletedAt:   "2022-03-06 13:53:00",
				Plot:        "2022-03-06 13:53:00",
				Released:    true,
			},
			mock:        movieStore.EXPECT().CreateMovie(gomock.Any(), gomock.Any()).Return(&model.MovieModel{}, nil),
			expectedErr: nil,
		},
	}
	h := NewMovieServiceHandler(movieStore)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, errEmp := h.InsertMovieService(ctx, &tc.movie)

		if errEmp != nil && errEmp.Error() != tc.expectedErr.Error() {
			t.Errorf("Got %v , Want : %v", tc.expectedErr, errEmp)
		}
	}
}

func TestService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieStore := store.NewMockMovieInterface(ctrl)

	tcs := []struct {
		mock        *gomock.Call
		expectedErr error
		ID          int
	}{
		{
			ID:          1,
			expectedErr: nil,
			mock:        movieStore.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil),
		},
		{
			ID:          -1,
			expectedErr: errors.InvalidParam{Param: []string{"id"}},
		},
		{
			ID:          1,
			expectedErr: errors.Error("prep"),
			mock:        movieStore.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.Error("prep")),
		},
	}
	h := NewMovieServiceHandler(movieStore)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())

		_, errEmp := h.GetByIDService(ctx, tc.ID)

		if errEmp != nil && errEmp.Error() != tc.expectedErr.Error() {
			t.Errorf("Got %v , Want : %v", tc.expectedErr, errEmp)
		}
	}
}

func TestService_DeleteByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieStore := store.NewMockMovieInterface(ctrl)
	tcs := []struct {
		mock        []*gomock.Call
		expectedErr error
		ID          int
	}{
		{
			ID:          -1,
			expectedErr: errors.Error("negative ID found"),
		},
		{
			ID:          1,
			mock:        []*gomock.Call{movieStore.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)},
			expectedErr: sql.ErrNoRows,
		},
		{
			ID: 1,
			mock: []*gomock.Call{
				movieStore.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, nil),
				movieStore.EXPECT().DeleteByID(gomock.Any(), gomock.Any()).Return(nil),
			},
			expectedErr: nil,
		},
	}
	h := NewMovieServiceHandler(movieStore)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		errEmp := h.DeleteByIDService(ctx, tc.ID)

		if errEmp != nil && errEmp.Error() != tc.expectedErr.Error() {
			t.Errorf("Got %v , Want : %v", tc.expectedErr, errEmp)
		}
	}
}

func TestService_UpdateByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieStore := store.NewMockMovieInterface(ctrl)
	tcs := []struct {
		mock        []*gomock.Call
		expectedErr error
		movie       model.MovieModel
	}{
		{
			movie: model.MovieModel{
				ID:          -1,
				Name:        "Richesh",
				Genre:       "comedy",
				Rating:      4.5,
				ReleaseDate: "2014-12-24",
				UpdatedAt:   "2022-03-06 19:29:33",
				CreatedAt:   "2022-03-06 13:53:00",
				DeletedAt:   "2022-03-06 13:53:00",
				Plot:        "2022-03-06 13:53:00",
				Released:    true,
			},

			expectedErr: errors.Error("negative ID found"),
		},
		{
			movie: model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "",
				Rating:      4.5,
				ReleaseDate: "2014-12-24",
				UpdatedAt:   "2022-03-06 19:29:33",
				CreatedAt:   "2022-03-06 13:53:00",
				DeletedAt:   "2022-03-06 13:53:00",
				Plot:        "2022-03-06 13:53:00",
				Released:    true,
			},
			mock: []*gomock.Call{
				movieStore.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows),
			},
			expectedErr: sql.ErrNoRows,
		},
		{
			movie: model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "comedy",
				Rating:      4.5,
				ReleaseDate: "2014-12-24",
				UpdatedAt:   "2022-03-06 19:29:33",
				CreatedAt:   "2022-03-06 13:53:00",
				DeletedAt:   "2022-03-06 13:53:00",
				Plot:        "2022-03-06 13:53:00",
				Released:    true,
			},
			mock: []*gomock.Call{
				movieStore.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&model.MovieModel{}, errors.Error("prep")),
				movieStore.EXPECT().UpdateByID(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows),
			},
			expectedErr: sql.ErrNoRows,
		},
		{
			movie: model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "comedy",
				Rating:      4.5,
				ReleaseDate: "2014-12-24",
				UpdatedAt:   "2022-03-06 19:29:33",
				CreatedAt:   "2022-03-06 13:53:00",
				DeletedAt:   "2022-03-06 13:53:00",
				Plot:        "2022-03-06 13:53:00",
				Released:    true,
			},
			mock: []*gomock.Call{
				movieStore.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(&model.MovieModel{}, errors.Error("prep")),
				movieStore.EXPECT().UpdateByID(gomock.Any(), gomock.Any()).Return(&model.MovieModel{}, nil),
			},
			expectedErr: sql.ErrNoRows,
		},
	}

	h := NewMovieServiceHandler(movieStore)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, errEmp := h.UpdatedByIDService(ctx, &tc.movie)

		if errEmp != nil && errEmp.Error() != tc.expectedErr.Error() {
			t.Errorf("Got %v , Want : %v", tc.expectedErr, errEmp)
		}
	}
}

func TestService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieStore := store.NewMockMovieInterface(ctrl)
	tcs := []struct {
		mock        *gomock.Call
		expectedErr error
	}{
		{
			expectedErr: sql.ErrNoRows,
			mock:        movieStore.EXPECT().GetAll(gomock.Any()).Return(nil, sql.ErrNoRows),
		},
		{
			expectedErr: nil,
			mock:        movieStore.EXPECT().GetAll(gomock.Any()).Return(&[]model.MovieModel{}, nil),
		},
	}

	h := NewMovieServiceHandler(movieStore)

	for _, tc := range tcs {
		ctx := gofr.NewContext(nil, nil, gofr.New())
		_, errEmp := h.GetAllService(ctx)

		if errEmp != nil && errEmp.Error() != tc.expectedErr.Error() {
			t.Errorf("Got %v , Want : %v", tc.expectedErr, errEmp)
		}
	}
}
