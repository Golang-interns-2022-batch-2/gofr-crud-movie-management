package movie

import (
	"bytes"

	"developer.zopsmart.com/go/gofr/pkg/errors"

	"encoding/json"
	"net/http"

	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/models"
	service "github.com/anushi/newbatch/gofr-crud-movie-management/gofr-test/services"
	"github.com/golang/mock/gomock"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := service.NewMockMovie(ctrl)

	app := gofr.New()

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testMovie := models.Movie{
		ID:          3,
		Name:        "Silicon valley",
		Genre:       "Comedy",
		Rating:      4.5,
		ReleaseDate: date,
		UpdatedAt:   dateTime,
		CreatedAt:   dateTime,
		Plot: "Richard, a programmer, " +
			"creates an app called the " +
			"Pied Piper and tries to get " +
			"investors for it. Meanwhile, " +
			"five other programmers struggle" +
			" to make their mark in " +
			"Silicon Valley.", Released: true}

	testCases := []struct {
		desc   string
		id     int
		experr error
		mock   []*gomock.Call
	}{
		{desc: "Get success",
			id:     3,
			experr: nil,
			mock: []*gomock.Call{
				mockService.EXPECT().GetByID(gomock.Any(), 3).Return(&testMovie, nil),
			},
		},
		{desc: "Get fail",
			id:     -1,
			experr: errors.InvalidParam{Param: []string{"id"}},
			mock: []*gomock.Call{
				mockService.EXPECT().GetByID(gomock.Any(), -1).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
			},
		},
	}

	s := New(mockService)

	for _, tc := range testCases {
		r := httptest.NewRequest("GET", "/movie/{id}"+strconv.Itoa(tc.id), nil)

		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": strconv.Itoa(tc.id),
		})

		_, err := s.GetByID(ctx)

		if !reflect.DeepEqual(err, tc.experr) {
			t.Errorf("Expected Error: %v, Got: %v", tc.experr, err)
		}
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockMovie(ctrl)
	app := gofr.New()
	s := New(mockService)

	testCases := []struct {
		id         int
		statusCode int
		err        error
		mock       *gomock.Call
	}{
		{id: 3, statusCode: http.StatusOK, err: nil, mock: mockService.EXPECT().Delete(gomock.Any(), 3).Return(nil)},
		{id: 7, statusCode: http.StatusBadRequest, err: errors.InvalidParam{Param: []string{"id"}},
			mock: mockService.EXPECT().Delete(gomock.Any(), 7).Return(errors.InvalidParam{Param: []string{"id"}})},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest("DELETE", "/movie/{id}"+strconv.Itoa(tc.id), nil)

		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": strconv.Itoa(tc.id),
		})

		_, err := s.Delete(ctx)

		if !reflect.DeepEqual(err, tc.err) {
			t.Errorf("Expected Err: %v, Got: %v", tc.err, err)
		}
	}
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockMovie(ctrl)
	app := gofr.New()
	s := New(mockService)

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testMovie := models.Movie{ID: 3, Name: "Silicon valley",
		Genre: "Comedy", Rating: 4.5,
		ReleaseDate: date, UpdatedAt: dateTime,
		CreatedAt: dateTime,
		Plot: "Richard, a programmer, " +
			"creates an app called the" +
			" Pied Piper and tries to get" +
			" investors for it. Meanwhile," +
			" five other programmers struggle" +
			" to make their mark in " +
			"Silicon Valley.", Released: true}

	testMovie1 := models.Movie{ID: 7, Name: "Silicon valley",
		Genre: "Comedy", Rating: 4.5,
		ReleaseDate: date, UpdatedAt: dateTime,
		CreatedAt: dateTime,
		Plot: "Richard, a programmer, " +
			"creates an app called the" +
			" Pied Piper and tries to get" +
			" investors for it. Meanwhile," +
			" five other programmers struggle" +
			" to make their mark in " +
			"Silicon Valley.", Released: true}

	tests := []struct {
		desc               string
		id                 int
		expectedStatusCode int
		body               models.Movie
		mockCall           *gomock.Call
		expErr             error
	}{
		{
			desc:               "Success",
			id:                 3,
			body:               testMovie,
			expectedStatusCode: http.StatusOK,
			mockCall:           mockService.EXPECT().Update(gomock.Any(), &testMovie).Return(&testMovie, nil),
			expErr:             nil,
		},
		{
			desc:               "Failure",
			id:                 7,
			body:               testMovie1,
			expectedStatusCode: http.StatusInternalServerError,
			mockCall:           mockService.EXPECT().Update(gomock.Any(), &testMovie1).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
			expErr:             errors.InvalidParam{Param: []string{"id"}},
		},
	}

	for _, tc := range tests {
		body, _ := json.Marshal(tc.body)

		r := httptest.NewRequest("PUT", "/movie/{id}"+strconv.Itoa(tc.id), bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		ctx.SetPathParams(map[string]string{
			"id": strconv.Itoa(tc.id),
		})

		_, err := s.Update(ctx)

		if !reflect.DeepEqual(tc.expErr, err) {
			t.Errorf("Expected Error: %v, Got: %v", tc.expErr, err)
		}
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockMovie(ctrl)
	app := gofr.New()
	s := New(mockService)

	date, _ := time.Parse(time.RFC3339, "2014-12-17")

	dateTime := time.Date(2014, 12, 17, 13, 39, 41, 0, time.UTC)

	testMovie := models.Movie{Name: "Silicon valley",
		Genre:       "Comedy",
		Rating:      4.5,
		ReleaseDate: date,
		UpdatedAt:   dateTime,
		CreatedAt:   dateTime,
		Plot: "Richard, a programmer," +
			" creates an app called the " +
			"Pied Piper and tries to get" +
			" investors for it. Meanwhile," +
			" five other programmers struggle" +
			" to make their mark in " +
			"Silicon Valley.", Released: true}

	tests := []struct {
		desc               string
		mov                models.Movie
		expectedStatusCode int
		expErr             error
		mockCall           *gomock.Call
	}{
		{
			desc:               "Success Case",
			mov:                testMovie,
			expectedStatusCode: http.StatusOK,
			expErr:             nil,
			mockCall:           mockService.EXPECT().Create(gomock.Any(), &testMovie).Return(&testMovie, nil),
		},
		{
			desc:               "Failure Case",
			mov:                models.Movie{},
			expectedStatusCode: http.StatusInternalServerError,
			expErr:             errors.InvalidParam{Param: []string{"id"}},
			mockCall:           mockService.EXPECT().Create(gomock.Any(), &models.Movie{}).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
		},
	}

	for _, test := range tests {
		body, _ := json.Marshal(test.mov)

		r := httptest.NewRequest("POST", "/movie", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := s.Create(ctx)

		if !reflect.DeepEqual(err, test.expErr) {
			t.Errorf("Expected Error: %v, Got: %v", test.expErr, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := service.NewMockMovie(ctrl)
	app := gofr.New()
	s := New(mockService)

	testCase := []struct {
		desc               string
		id                 int
		expectedStatusCode int
		expErr             error
		mockCall           []*gomock.Call
	}{
		{
			desc:               "Get All Success",
			id:                 3,
			expectedStatusCode: http.StatusOK,
			expErr:             nil,
			mockCall: []*gomock.Call{
				mockService.EXPECT().GetAll(gomock.Any()).Return([]*models.Movie{}, nil),
			},
		},
		{
			desc:               "Get All Fail",
			id:                 0,
			expectedStatusCode: http.StatusBadRequest,
			expErr:             errors.InvalidParam{Param: []string{"id"}},
			mockCall: []*gomock.Call{
				mockService.EXPECT().GetAll(gomock.Any()).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
			},
		},
	}

	for _, test := range testCase {
		r := httptest.NewRequest("GET", "/movie", nil)
		w := httptest.NewRecorder()

		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)
		ctx := gofr.NewContext(res, req, app)

		_, err := s.GetAll(ctx)
		if !reflect.DeepEqual(err, test.expErr) {
			t.Errorf("Expected Error: %v, Got: %v", test.expErr, err)
		}
	}
}
