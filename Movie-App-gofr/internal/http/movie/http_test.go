package movie

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/model"
	"github.com/RicheshZopsmart/Movie-App-gofr/internal/service"

	"github.com/golang/mock/gomock"
)

func TestCreateMovieRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configServ := service.NewMockMovieServiceInterface(ctrl)

	tcs := []struct {
		Body       []byte
		statusCode int64
		err        error
		mock       *gomock.Call
	}{
		{
			Body: []byte(`{{
				"name": "Tilicon valley",
				"genre": "Comedy Central",
				"rating": 1.5,
				"plot": "Richard ",
				"released": true,
				"releaseDate":"2014-12-17"
			}`),
			statusCode: 400,
			err:        errors.InvalidParam{Param: []string{"body"}},
		},
		{
			Body: []byte(`{
				"name": "Tilicon valley",
				"genre": "Comedy Central",
				"rating": 1.5,
				"plot": "Richard ",
				"released": true,
				"releaseDate":"2014-12-17"
			}`),
			statusCode: 200,
			err:        nil,
			mock:       configServ.EXPECT().InsertMovieService(gomock.Any(), gomock.Any()).Return(nil, nil),
		},
	}

	h := New(configServ)

	for _, tc := range tcs {
		req := httptest.NewRequest(http.MethodPost, "http://localhost:8090/movies", bytes.NewReader(tc.Body))
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		_, err := h.CreateMovieRequest(ctx)

		if err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Want : %v, Got : %v", tc.err.Error(), err.Error())
		}

	}
}
func TestDeleteByIdRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configServ := service.NewMockMovieServiceInterface(ctrl)

	tcs := []struct {
		id         string
		statusCode int64
		err        error
		mock       *gomock.Call
	}{
		{
			id:         "1",
			statusCode: 200,
			err:        nil,
			mock:       configServ.EXPECT().DeleteByIDService(gomock.Any(), gomock.Any()).Return(nil),
		},
		{
			id:         "1",
			statusCode: 404,
			err:        errors.EntityNotFound{Entity: "movie", ID: "1"},
			mock:       configServ.EXPECT().DeleteByIDService(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows),
		},
		{
			id:         "1",
			statusCode: 500,
			err:        errors.Error("internal server error"),
			mock:       configServ.EXPECT().DeleteByIDService(gomock.Any(), gomock.Any()).Return(errors.Error("internal server error")),
		},
		{
			id:         "abcd",
			statusCode: 400,
			err:        errors.InvalidParam{Param: []string{"id"}},
		},
	}

	h := New(configServ)

	for _, tc := range tcs {

		req := httptest.NewRequest(http.MethodDelete, "http://localhost:8090/movies/1", nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		ctx.SetPathParams(map[string]string{"id": tc.id})
		_, err := h.DeleteByIDRequest(ctx)

		if err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Want : %v, Got : %v", tc.err.Error(), err.Error())
		}

	}
}

func TestUpdateByIdRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configServ := service.NewMockMovieServiceInterface(ctrl)

	mObj := &model.MovieModel{
		ID:          1,
		Name:        "Silicon Valley",
		Genre:       "comedy",
		Rating:      4.5,
		Plot:        "Richard",
		Released:    true,
		ReleaseDate: "2014-12-17",
	}

	JSONObj, _ := json.Marshal(mObj)

	tcs := []struct {
		id         string
		Body       []byte
		statusCode int
		movieObj   *model.MovieModel
		mock       *gomock.Call
		err        error
	}{
		// Failure Case
		{
			Body:       JSONObj,
			id:         "1",
			statusCode: 500,
			movieObj: &model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "comedy",
				Rating:      4.5,
				Plot:        "Richard",
				Released:    true,
				ReleaseDate: "2014-12-17",
			},
			err:  errors.Error("internal server error"),
			mock: configServ.EXPECT().UpdatedByIDService(gomock.Any(), mObj).Return(mObj, errors.Error("internal server error")),
		},
		// 404 Error
		{
			Body:       JSONObj,
			id:         "1",
			statusCode: 500,
			movieObj: &model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "comedy",
				Rating:      4.5,
				Plot:        "Richard",
				Released:    true,
				ReleaseDate: "2014-12-17",
			},
			err:  errors.EntityNotFound{Entity: "movie", ID: "1"},
			mock: configServ.EXPECT().UpdatedByIDService(gomock.Any(), mObj).Return(mObj, sql.ErrNoRows),
		},
		// Success Case
		{
			Body:       JSONObj,
			id:         "1",
			statusCode: 200,
			movieObj: &model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "comedy",
				Rating:      4.5,
				Plot:        "Richard",
				Released:    true,
				ReleaseDate: "2014-12-17",
			},
			err:  nil,
			mock: configServ.EXPECT().UpdatedByIDService(gomock.Any(), mObj).Return(mObj, nil),
		},
		{
			Body:       JSONObj,
			id:         "abcd",
			statusCode: 200,
			movieObj: &model.MovieModel{
				ID:          1,
				Name:        "Silicon Valley",
				Genre:       "comedy",
				Rating:      4.5,
				Plot:        "Richard",
				Released:    true,
				ReleaseDate: "2014-12-17",
			},
			err: errors.InvalidParam{Param: []string{"id"}},
		},
	}

	h := New(configServ)

	for _, tc := range tcs {

		req := httptest.NewRequest(http.MethodDelete, "http://localhost:8090/movies/1", bytes.NewReader(tc.Body))
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)
		ctx.SetPathParams(map[string]string{"id": tc.id})
		// Mocking
		_, err := h.UpdateByIDRequest(ctx)

		if err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Got : %v want : %v", err.Error(), tc.err.Error())
		}
	}
}
func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configServ := service.NewMockMovieServiceInterface(ctrl)
	tcs := []struct {
		statusCode int64
		err        error
		mock       *gomock.Call
	}{
		// GetAll() error mock
		{
			statusCode: 500,
			err:        errors.Error("internal server error"),
			mock:       configServ.EXPECT().GetAllService(gomock.Any()).Return(nil, errors.Error("internal server error")),
		},
		{
			statusCode: 200,
			err:        nil,
			mock:       configServ.EXPECT().GetAllService(gomock.Any()).Return(nil, nil),
		},
	}

	for _, tc := range tcs {

		h := New(configServ)
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8090/movies", nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		_, err := h.GetAllRequest(ctx)

		if err != nil && err.Error() != tc.err.Error() {
			t.Errorf("Got : %v want : %v", err.Error(), tc.err.Error())
		}
	}
}
func TestGetById(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	configServ := service.NewMockMovieServiceInterface(ctrl)

	tcs := []struct {
		ID          string
		expectedErr error
		mockQ       *gomock.Call
	}{
		{
			ID:          "1",
			expectedErr: nil,
			mockQ: configServ.EXPECT().GetByIDService(gomock.Any(), 1).Return(
				&model.MovieModel{}, nil),
		},

		{
			ID:          "abcd",
			expectedErr: errors.InvalidParam{Param: []string{"id"}},
		},
	}
	h := New(configServ)

	for _, tc := range tcs {

		req := httptest.NewRequest(http.MethodGet, "http://localhost:8090/movies/1", nil)
		ctx := gofr.NewContext(responder.NewContextualResponder(httptest.NewRecorder(), req), request.NewHTTPRequest(req), nil)

		ctx.SetPathParams(map[string]string{"id": tc.ID})

		_, err := h.GetByIDRequest(ctx)

		if !reflect.DeepEqual(tc.expectedErr, err) {
			t.Errorf("Failed Tc Err : %v , Expected Err : %v", tc.expectedErr.Error(), err.Error())
		}
	}
}
