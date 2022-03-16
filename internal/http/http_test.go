package http

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	gofrerr "developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/golang/mock/gomock"
	"github.com/shivam/Crud_Gofr/internal/models"
	"github.com/shivam/Crud_Gofr/internal/service"
)

func TestGetById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MovieServiceMock := service.NewMockInterface(mockCtrl)

	mov := models.Movie{
		ID:          1,
		Name:        "MazeRunner",
		Genre:       "Action",
		Rating:      5.0,
		ReleaseDate: "2022-12-17",
		UpdatedAt:   "",
		CreatedAt:   "",
		Plot:        "PLOT!",
		Released:    true,
	}

	tc := []struct {
		id     string
		desc   string
		output interface{}
		err    error
		mock   *gomock.Call
	}{
		{
			id:     "1",
			desc:   "success case",
			output: &mov,
			mock:   MovieServiceMock.EXPECT().GetByIDService(gomock.Any(), "1").Return(&mov, nil),
		},

		{
			id:     "1",
			desc:   "Error case",
			output: nil,
			err: gofrerr.EntityNotFound{
				Entity: "Movie",
				ID:     "1",
			},
			mock: MovieServiceMock.EXPECT().GetByIDService(gomock.Any(), "1").Return(nil, gofrerr.EntityNotFound{
				Entity: "Movie",
				ID:     "1",
			}),
		},
	}

	for _, tc := range tc {
		tc := tc

		t.Run("", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/movies", nil)
			w := httptest.NewRecorder()
			req := request.NewHTTPRequest(r)
			res := responder.NewContextualResponder(w, r)

			ctx := gofr.NewContext(res, req, nil)

			ctx.SetPathParams(map[string]string{
				"id": tc.id,
			})

			httpHandler := New(MovieServiceMock)
			response, _ := httpHandler.GetByID(ctx)

			if !reflect.DeepEqual(response, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, response)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MovieServiceMock := service.NewMockInterface(mockCtrl)

	mov := []*models.Movie{
		{
			ID:          1,
			Name:        "MazeRunner",
			Genre:       "Action",
			Rating:      5,
			ReleaseDate: "2022-12-17",
			UpdatedAt:   "",
			CreatedAt:   "",
			Plot:        "PLOT!",
			Released:    true,
		},
	}

	tcs := []struct {
		desc   string
		err    error
		output interface{}
		mock   *gomock.Call
	}{
		{
			desc: "success case",

			output: response{
				Data: data{mov},
			},
			mock: MovieServiceMock.EXPECT().GetAllService(gomock.Any()).Return(mov, nil),
		},
		{
			desc:   "internal server error",
			output: nil,
			err:    gofrerr.Error("Internal Server Error"),
			mock:   MovieServiceMock.EXPECT().GetAllService(gomock.Any()).Return(nil, gofrerr.Error("Internal Server Error")),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/movies", nil)
			w := httptest.NewRecorder()
			req := request.NewHTTPRequest(r)
			res := responder.NewContextualResponder(w, r)

			ctx := gofr.NewContext(res, req, nil)

			httpHandler := New(MovieServiceMock)
			response, _ := httpHandler.GetAll(ctx)

			if !reflect.DeepEqual(response, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, response)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MovieServiceMock := service.NewMockInterface(mockCtrl)

	tcs := []struct {
		desc   string
		id     string
		err    error
		output interface{}
		mock   *gomock.Call
	}{
		{
			desc:   "success case",
			id:     "1",
			output: "Deleted successfully",
			err:    nil,
			mock:   MovieServiceMock.EXPECT().DeleteService(gomock.Any(), "1").Return(nil),
		},

		{
			desc:   "Fail Case",
			id:     "-2",
			output: nil,
			err:    gofrerr.Error("Internal Servor Error"),
			mock:   MovieServiceMock.EXPECT().DeleteService(gomock.Any(), "-2").Return(gofrerr.Error("Internal Servor Error")),
		},
	}

	for _, testcase := range tcs {
		testcase := testcase

		r := httptest.NewRequest("GET", "/movies", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		ctx.SetPathParams(map[string]string{
			"id": testcase.id,
		})

		httpHandler := New(MovieServiceMock)
		response, _ := httpHandler.DeleteByID(ctx)

		if !reflect.DeepEqual(response, testcase.output) {
			t.Errorf("Expected %v got %v", testcase.output, response)
		}
	}
}

func TestPost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MovieServiceMock := service.NewMockInterface(mockCtrl)

	movie := models.Movie{
		ID:          1,
		Name:        "MazeRunner",
		Genre:       "Action",
		Rating:      5,
		ReleaseDate: "2022-12-17",
		Plot:        "PLOT!",
		Released:    true,
	}

	requestBody, _ := json.Marshal(movie)

	tcs := []struct {
		desc   string
		body   []byte
		output interface{}
		err    error
		mock   *gomock.Call
	}{

		{
			desc:   "success case",
			body:   requestBody,
			err:    nil,
			output: &movie,
			mock:   MovieServiceMock.EXPECT().InsertService(gomock.Any(), &movie).Return(&movie, nil),
		},

		{
			desc:   "Bind Error",
			body:   []byte(`{"id":1"name":"MazeRunner","genre":"Action","rating":5.0,"releasedate": "2022-12-17","plot":"PLOT!","released":true}`),
			output: nil,
			err:    gofrerr.InvalidParam{Param: []string{"body"}},
		},

		{
			desc:   "Internal server error",
			body:   requestBody,
			output: nil,
			err:    gofrerr.Error("Internal server Error"),
			mock:   MovieServiceMock.EXPECT().InsertService(gomock.Any(), &movie).Return(nil, gofrerr.Error("Internal server Error")),
		},
	}

	for _, tc := range tcs {
		tc := tc

		t.Run("", func(t *testing.T) {
			r := httptest.NewRequest("POST", "/movies", bytes.NewReader(tc.body))
			w := httptest.NewRecorder()
			req := request.NewHTTPRequest(r)
			res := responder.NewContextualResponder(w, r)

			ctx := gofr.NewContext(res, req, nil)

			httpHandler := New(MovieServiceMock)
			response, _ := httpHandler.Create(ctx)

			if !reflect.DeepEqual(response, tc.output) {
				t.Errorf("Expected %v got %v", tc.output, response)
			}
		})
	}
}

func TestPut(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	MovieServiceMock := service.NewMockInterface(mockCtrl)

	movie := models.Movie{
		ID:          1,
		Name:        "MazeRunner",
		Genre:       "Action",
		Rating:      5,
		ReleaseDate: "2022-12-17",
		UpdatedAt:   "",
		CreatedAt:   "",
		Plot:        "PLOT!",
		Released:    true,
	}

	requestBody, _ := json.Marshal(movie)

	tcs := []struct {
		desc   string
		id     string
		body   []byte
		output interface{}
		err    error
		mock   *gomock.Call
	}{

		{
			desc: "success case",
			id:   "1",
			body: requestBody,

			output: &movie,
			mock: MovieServiceMock.EXPECT().UpdatedService(gomock.Any(), &movie, "1").
				Return(&models.Movie{
					ID:          1,
					Name:        "MazeRunner",
					Genre:       "Action",
					Rating:      5,
					ReleaseDate: "2022-12-17",
					Plot:        "PLOT!",
					Released:    true,
				}, nil),
		},

		{
			id:   "1",
			desc: "Bind Error",
			body: []byte(`{"id":"1", "name":"MazeRunner" genre":"Action","rating":5.0, ` +
				`"releasedate": "2022-12-17","plot":"PLOT!","released":true}`),
			output: nil,
			err:    gofrerr.InvalidParam{Param: []string{"body"}},
		},

		{
			desc:   "Internal server error",
			id:     "1",
			body:   requestBody,
			output: nil,
			err:    gofrerr.Error("Internal Server Error"),
			mock:   MovieServiceMock.EXPECT().UpdatedService(gomock.Any(), &movie, "1").Return(nil, gofrerr.Error("Internal Server Error")),
		},
	}

	for _, tc := range tcs {
		tc := tc

		r := httptest.NewRequest("GET", "/movies", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		httpHandler := New(MovieServiceMock)
		response, _ := httpHandler.UpdateByID(ctx)

		if !reflect.DeepEqual(response, tc.output) {
			t.Errorf("Expected %v got %v", tc.output, response)
		}
	}
}
