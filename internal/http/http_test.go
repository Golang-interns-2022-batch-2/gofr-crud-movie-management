package http

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/golang/mock/gomock"
	"github.com/iamkakashi/movie-gofr/internal/model"
	"github.com/iamkakashi/movie-gofr/internal/service"
	"gopkg.in/guregu/null.v4"
)

func TestGetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	parsedTime, _ := time.Parse(time.RFC3339, "2022-03-05T18:44:05Z")
	mov := model.Movie{
		ID:          4,
		Name:        null.StringFrom("Deadpool"),
		Genre:       null.StringFrom("Comedy, Action"),
		Rating:      null.FloatFrom(4.8),
		ReleaseDate: null.StringFrom("2014-12-17"),
		UpdatedAt:   null.TimeFrom(parsedTime),
		CreatedAt:   null.TimeFrom(parsedTime),
		Plot:        null.StringFrom("This is a superhero movie"),
		Released:    null.BoolFrom(true),
	}
	MovieServiceMock := service.NewMockMovieServicer(mockCtrl)
	tcs := []struct {
		id     string
		desc   string
		output interface{}
		mock   *gomock.Call
		erOut  error
	}{
		{
			desc: "success case",
			id:   "4",
			output: apiResponse{
				Code:   200,
				Status: "SUCCESS",
				Data:   data{&mov},
			},
			erOut: nil,
			mock:  MovieServiceMock.EXPECT().GetByID(gomock.Any(), 4).Return(&mov, nil),
		},
		{
			desc:   "id doesn't exist",
			id:     "4",
			output: nil,
			erOut: errors.EntityNotFound{
				Entity: "movie",
				ID:     "4",
			},
			mock: MovieServiceMock.EXPECT().GetByID(gomock.Any(), 4).Return(nil, errors.EntityNotFound{
				Entity: "movie",
				ID:     "4",
			}),
		},
		{
			desc:   "invalid id",
			id:     "3b",
			output: nil,
			erOut:  errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc:   "missing id",
			id:     "",
			output: nil,
			erOut:  errors.MissingParam{Param: []string{"id"}},
		},
	}

	for _, tc := range tcs {
		tc := tc

		r := httptest.NewRequest("GET", "/movies", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		httpHandler := New(MovieServiceMock)
		resp, err := httpHandler.GetByID(ctx)

		if !reflect.DeepEqual(tc.erOut, err) {
			t.Errorf("Expected %v got %v", tc.erOut, err)
		}

		if !reflect.DeepEqual(tc.output, resp) {
			t.Errorf("Expected %v got %v", tc.output, resp)
		}
	}
}

func TestGet(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	parsedTime, _ := time.Parse(time.RFC3339, "2022-03-05T18:44:05Z")

	mov := []*model.Movie{
		{
			ID:          4,
			Name:        null.StringFrom("Deadpool"),
			Genre:       null.StringFrom("Comedy, Action"),
			Rating:      null.FloatFrom(4.8),
			ReleaseDate: null.StringFrom("2014-12-17"),
			UpdatedAt:   null.TimeFrom(parsedTime),
			CreatedAt:   null.TimeFrom(parsedTime),
			Plot:        null.StringFrom("This is a superhero movie"),
			Released:    null.BoolFrom(true),
		},
	}
	MovieServiceMock := service.NewMockMovieServicer(mockCtrl)

	tcs := []struct {
		desc   string
		output interface{}
		erOut  error
		mock   *gomock.Call
	}{
		{
			desc: "success case",
			output: apiResponse{
				Code:   200,
				Status: "SUCCESS",
				Data:   data{mov},
			},
			erOut: nil,
			mock:  MovieServiceMock.EXPECT().Get(gomock.Any()).Return(mov, nil),
		},
		{
			desc: "failure",

			output: nil,
			erOut:  errors.Error("no movie found"),
			mock:   MovieServiceMock.EXPECT().Get(gomock.Any()).Return(nil, errors.Error("no movie found")),
		},
	}

	for _, tc := range tcs {
		tc := tc

		r := httptest.NewRequest("GET", "/movies", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		httpHandler := New(MovieServiceMock)
		resp, err := httpHandler.Get(ctx)

		if !reflect.DeepEqual(tc.erOut, err) {
			t.Errorf("Expected %v got %v", tc.erOut, err)
		}

		if !reflect.DeepEqual(tc.output, resp) {
			t.Errorf("Expected %v got %v", tc.output, resp)
		}
	}
}

func TestPost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	movie := model.Movie{
		ID:          4,
		Name:        null.StringFrom("Deadpool"),
		Genre:       null.StringFrom("Comedy, Action"),
		Rating:      null.FloatFrom(4.8),
		ReleaseDate: null.StringFrom("2023-12-17"),
		Plot:        null.StringFrom("This is a superhero movie"),
		Released:    null.BoolFrom(false),
	}
	body := model.Movie{
		Name:        null.StringFrom("Deadpool"),
		Genre:       null.StringFrom("Comedy, Action"),
		Rating:      null.FloatFrom(4.8),
		ReleaseDate: null.StringFrom("2023-12-17"),
		Plot:        null.StringFrom("This is a superhero movie"),
		Released:    null.BoolFrom(false),
	}
	requestBody, _ := json.Marshal(body)
	MovieServiceMock := service.NewMockMovieServicer(mockCtrl)

	tcs := []struct {
		desc   string
		body   []byte
		output interface{}
		erOut  error
		mock   *gomock.Call
	}{

		{
			desc: "success case",

			body: requestBody,

			output: apiResponse{
				Code:   200,
				Status: "SUCCESS",
				Data:   data{&movie},
			},
			erOut: nil,
			mock: MovieServiceMock.EXPECT().Create(gomock.Any(), &body).
				Return(&movie, nil),
		},
		{
			desc: "binding error",

			body: []byte(`{"name":"Deadpool"
			"genre":"Comedy, Action",
			"rating":4.8,
			"releasedate": "2023-12-17",
			"plot":"This is a superhero movie",
			"released":false}`),
			output: nil,
			erOut:  errors.InvalidParam{Param: []string{"body"}},
		},
		{
			desc:   "internal server error",
			body:   requestBody,
			output: nil,
			erOut:  errors.Error("internal server error"),
			mock:   MovieServiceMock.EXPECT().Create(gomock.Any(), &body).Return(nil, errors.Error("internal server error")),
		},
	}

	for _, tc := range tcs {
		tc := tc
		r := httptest.NewRequest("POST", "/movies", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		httpHandler := New(MovieServiceMock)
		resp, err := httpHandler.Post(ctx)

		if !reflect.DeepEqual(tc.erOut, err) {
			t.Errorf("Expected %v got %v", tc.erOut, err)
		}

		if !reflect.DeepEqual(tc.output, resp) {
			t.Errorf("Expected %v got %v", tc.output, resp)
		}
	}
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MovieServiceMock := service.NewMockMovieServicer(mockCtrl)

	tcs := []struct {
		desc   string
		id     string
		output interface{}
		erOut  error
		mock   *gomock.Call
	}{
		{
			desc:   "success case",
			id:     "4",
			output: "Deleted successfully",
			erOut:  nil,
			mock:   MovieServiceMock.EXPECT().Delete(gomock.Any(), 4).Return(nil),
		},
		{
			desc:   "missing id",
			id:     "",
			output: nil,
			erOut:  errors.MissingParam{Param: []string{"id"}},
		},
		{
			desc:   "invalid id",
			id:     "3b",
			output: nil,
			erOut:  errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc:   "internal server error",
			id:     "4",
			output: nil,
			erOut:  errors.Error("internal server error"),
			mock:   MovieServiceMock.EXPECT().Delete(gomock.Any(), 4).Return(errors.Error("internal server error")),
		},
	}

	for _, tc := range tcs {
		r := httptest.NewRequest("DELETE", "/movies", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		httpHandler := New(MovieServiceMock)
		resp, err := httpHandler.Delete(ctx)

		if !reflect.DeepEqual(tc.erOut, err) {
			t.Errorf("Expected %v got %v", tc.erOut, err)
		}

		if !reflect.DeepEqual(tc.output, resp) {
			t.Errorf("Expected %v got %v", tc.output, resp)
		}
	}
}

func TestPut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	movie := model.Movie{
		ID:          4,
		Name:        null.StringFrom("Deadpool"),
		Genre:       null.StringFrom("Comedy, Action"),
		Rating:      null.FloatFrom(4.8),
		ReleaseDate: null.StringFrom("2023-12-17"),
		Plot:        null.StringFrom("This is a superhero movie"),
		Released:    null.BoolFrom(false),
	}
	body := model.Movie{
		Name:  null.StringFrom("Deadpool"),
		Genre: null.StringFrom("Comedy, Action"),
	}
	requestBody, _ := json.Marshal(body)

	MovieServiceMock := service.NewMockMovieServicer(mockCtrl)

	tcs := []struct {
		desc   string
		id     string
		body   []byte
		output interface{}
		erOut  error
		mock   *gomock.Call
	}{
		{
			desc: "binding error",
			id:   "4",
			body: []byte(
				`{"name":"Deadpool"
				"genre":"Comedy, Action",
				"rating":4.8,
				"releasedate": "2023-12-17",
				"plot":"This is a superhero movie",
				"released":false}`,
			),
			output: nil,
			erOut:  errors.InvalidParam{Param: []string{"body"}},
		},
		{
			desc: "success case",
			id:   "4",
			body: requestBody,

			output: apiResponse{
				Code:   200,
				Status: "SUCCESS",
				Data:   data{&movie},
			},
			erOut: nil,
			mock: MovieServiceMock.EXPECT().Update(gomock.Any(), &model.Movie{
				ID:    4,
				Name:  null.StringFrom("Deadpool"),
				Genre: null.StringFrom("Comedy, Action"),
			}).
				Return(&model.Movie{
					ID:          4,
					Name:        null.StringFrom("Deadpool"),
					Genre:       null.StringFrom("Comedy, Action"),
					Rating:      null.FloatFrom(4.8),
					ReleaseDate: null.StringFrom("2023-12-17"),
					Plot:        null.StringFrom("This is a superhero movie"),
					Released:    null.BoolFrom(false),
				}, nil),
		},
		{
			desc:   "missing id",
			id:     "",
			output: nil,
			erOut:  errors.MissingParam{Param: []string{"id"}},
		},
		{
			desc:   "invalid id",
			id:     "3b",
			output: nil,
			erOut:  errors.InvalidParam{Param: []string{"id"}},
		},
		{
			desc: "fail case",
			id:   "4",
			body: requestBody,

			output: nil,
			erOut:  errors.Error("failed to update"),
			mock: MovieServiceMock.EXPECT().Update(gomock.Any(), &model.Movie{
				ID:    4,
				Name:  null.StringFrom("Deadpool"),
				Genre: null.StringFrom("Comedy, Action"),
			}).
				Return(nil, errors.Error("failed to update")),
		},
	}

	for _, tc := range tcs {
		tc := tc
		r := httptest.NewRequest("PUT", "/movies", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		httpHandler := New(MovieServiceMock)
		resp, err := httpHandler.Put(ctx)

		if !reflect.DeepEqual(tc.erOut, err) {
			t.Errorf("Expected %v got %v", tc.erOut, err)
		}

		if !reflect.DeepEqual(tc.output, resp) {
			t.Errorf("Expected %v got %v", tc.output, resp)
		}
	}
}
