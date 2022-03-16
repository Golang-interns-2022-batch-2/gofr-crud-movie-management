package movie

import (
	"bytes"
	"developer.zopsmart.com/go/gofr/pkg/errors"
	"developer.zopsmart.com/go/gofr/pkg/gofr"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"golangprog/gofr-test/Internal/models"
	"golangprog/gofr-test/Internal/service"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetByID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := service.NewMockMovie(mockCtrl)
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
		Plot: `"Richard, a programmer,
		creates an app called the Pied Piper and tries to get investors for it.
		Meanwhile, five other programmers struggle to make their mark in Silicon Valley."`,
		Released: true,
	}

	testCases := []struct {
		desc               string
		id                 string
		expectedStatusCode int
		expectederr        error
		mock               []*gomock.Call
	}{
		{desc: "Get success",
			id:                 "1",
			expectedStatusCode: 200,
			expectederr:        nil,
			mock: []*gomock.Call{
				mockService.EXPECT().GetByID(gomock.Any(), 1).Return(&testMovie, nil),
			},
		},
		{desc: "Get failure",
			id:                 "0",
			expectedStatusCode: 400,
			expectederr:        errors.InvalidParam{Param: []string{"id"}},
			mock: []*gomock.Call{
				mockService.EXPECT().GetByID(gomock.Any(), 0).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		r := httptest.NewRequest("GET", "/movies", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		httpHandler := New(mockService)
		_, err := httpHandler.GetByID(ctx)

		if !reflect.DeepEqual(tc.expectederr, err) {
			t.Errorf("Expected %v got %v", tc.expectederr, err)
		}
	}
}

func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := service.NewMockMovie(mockCtrl)
	testCases := []struct {
		desc               string
		id                 string
		expectedStatusCode int
		expectederr        error
		mock               []*gomock.Call
	}{
		{desc: "Get success",
			id:                 "1",
			expectedStatusCode: 200,
			expectederr:        nil,
			mock: []*gomock.Call{
				mockService.EXPECT().GetAll(gomock.Any()).Return([]*models.Movie{}, nil),
			},
		},
		{desc: "Get failure",
			id:                 "0",
			expectedStatusCode: 400,
			expectederr:        errors.Error("bad request"),
			mock: []*gomock.Call{
				mockService.EXPECT().GetAll(gomock.Any()).Return(nil, errors.Error("bad request")),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		r := httptest.NewRequest("GET", "/movies", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		httpHandler := New(mockService)
		_, err := httpHandler.GetAll(ctx)

		if !reflect.DeepEqual(tc.expectederr, err) {
			t.Errorf("Expected %v got %v", tc.expectederr, err)
		}
	}
}

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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
	requestBody1, _ := json.Marshal(testMovie)
	testMovie1 := models.Movie{
		ID:          0,
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
	requestBody2, _ := json.Marshal(testMovie1)
	mockService := service.NewMockMovie(mockCtrl)
	testCases := []struct {
		desc               string
		id                 string
		expectedStatusCode int
		body               []byte
		expectederr        error
		mock               *gomock.Call
	}{
		{desc: "Get success",
			id:                 "1",
			expectedStatusCode: 200,
			body:               requestBody1,
			expectederr:        nil,
			mock:               mockService.EXPECT().Create(gomock.Any(), &testMovie).Return(&testMovie, nil),
		},
		{desc: "Get failure",
			id:                 "0",
			expectedStatusCode: 400,
			body:               requestBody2,
			expectederr:        errors.InvalidParam{Param: []string{"id"}},
			mock:               mockService.EXPECT().Create(gomock.Any(), &testMovie1).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		r := httptest.NewRequest("POST", "/movies", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)

		httpHandler := New(mockService)
		_, err := httpHandler.Create(ctx)

		if !reflect.DeepEqual(tc.expectederr, err) {
			t.Errorf("Expected %v got %v", tc.expectederr, err)
		}
	}
}

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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
	requestBody1, _ := json.Marshal(testMovie)
	testMovie1 := models.Movie{
		ID:          0,
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
	requestBody2, _ := json.Marshal(testMovie1)
	mockService := service.NewMockMovie(mockCtrl)
	testCases := []struct {
		desc               string
		id                 string
		expectedStatusCode int
		body               []byte
		expectederr        error
		mock               *gomock.Call
	}{
		{desc: "Get success",
			id:                 "1",
			expectedStatusCode: 200,
			expectederr:        nil,
			body:               requestBody1,
			mock:               mockService.EXPECT().Update(gomock.Any(), &testMovie).Return(&testMovie, nil),
		},
		{desc: "Get failure",
			id:                 "0",
			expectedStatusCode: 400,
			expectederr:        errors.InvalidParam{Param: []string{"id"}},
			body:               requestBody2,
			mock:               mockService.EXPECT().Update(gomock.Any(), &testMovie1).Return(nil, errors.InvalidParam{Param: []string{"id"}}),
		},
	}

	for _, tc := range testCases {
		tc := tc
		r := httptest.NewRequest("PUT", "/movies", bytes.NewReader(tc.body))
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		httpHandler := New(mockService)
		_, err := httpHandler.Update(ctx)

		if !reflect.DeepEqual(tc.expectederr, err) {
			t.Errorf("Expected %v got %v", tc.expectederr, err)
		}
	}
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockService := service.NewMockMovie(mockCtrl)
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

	testCases := []struct {
		desc               string
		id                 string
		movieObj           models.Movie
		expectedStatusCode int
		expectederr        error
		mock               []*gomock.Call
	}{
		{desc: "Get success",
			id:                 "1",
			expectedStatusCode: 200,
			movieObj:           testMovie,
			expectederr:        nil,
			mock: []*gomock.Call{
				mockService.EXPECT().Delete(gomock.Any(), testMovie.ID).Return(nil),
			},
		},
		{desc: "Get failure",
			id:                 "0",
			expectedStatusCode: 400,
			movieObj:           testMovie1,
			expectederr:        errors.Error("bad request"),
			mock: []*gomock.Call{
				mockService.EXPECT().Delete(gomock.Any(), testMovie1.ID).Return(errors.Error("bad request")),
			},
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest("DELETE", "/movies", nil)
		w := httptest.NewRecorder()
		req := request.NewHTTPRequest(r)
		res := responder.NewContextualResponder(w, r)

		ctx := gofr.NewContext(res, req, nil)
		ctx.SetPathParams(map[string]string{
			"id": tc.id,
		})

		httpHandler := New(mockService)
		_, err := httpHandler.Delete(ctx)

		if !reflect.DeepEqual(tc.expectederr, err) {
			t.Errorf("Expected %v got %v", tc.expectederr, err)
		}
	}
}
