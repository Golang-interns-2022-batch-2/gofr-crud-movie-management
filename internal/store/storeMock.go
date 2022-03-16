// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/shivam/Crud_Gofr/internal/store (interfaces: MovieStorer)

// Package store is a generated GoMock package.
package store

import (
	gofr "developer.zopsmart.com/go/gofr/pkg/gofr"
	gomock "github.com/golang/mock/gomock"
	models "github.com/shivam/Crud_Gofr/internal/models"
	reflect "reflect"
)

// MockMovieStorer is a mock of MovieStorer interface
type MockMovieStorer struct {
	ctrl     *gomock.Controller
	recorder *MockMovieStorerMockRecorder
}

// MockMovieStorerMockRecorder is the mock recorder for MockMovieStorer
type MockMovieStorerMockRecorder struct {
	mock *MockMovieStorer
}

// NewMockMovieStorer creates a new mock instance
func NewMockMovieStorer(ctrl *gomock.Controller) *MockMovieStorer {
	mock := &MockMovieStorer{ctrl: ctrl}
	mock.recorder = &MockMovieStorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMovieStorer) EXPECT() *MockMovieStorerMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockMovieStorer) Create(arg0 *gofr.Context, arg1 *models.Movie) (*models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockMovieStorerMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMovieStorer)(nil).Create), arg0, arg1)
}

// DeleteByID mocks base method
func (m *MockMovieStorer) DeleteByID(arg0 *gofr.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID
func (mr *MockMovieStorerMockRecorder) DeleteByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockMovieStorer)(nil).DeleteByID), arg0, arg1)
}

// GetAll mocks base method
func (m *MockMovieStorer) GetAll(arg0 *gofr.Context) ([]*models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]*models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockMovieStorerMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockMovieStorer)(nil).GetAll), arg0)
}

// GetByID mocks base method
func (m *MockMovieStorer) GetByID(arg0 *gofr.Context, arg1 int) (*models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockMovieStorerMockRecorder) GetByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockMovieStorer)(nil).GetByID), arg0, arg1)
}

// Update mocks base method
func (m *MockMovieStorer) Update(arg0 *gofr.Context, arg1 int, arg2 *models.Movie) (*models.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockMovieStorerMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMovieStorer)(nil).Update), arg0, arg1, arg2)
}
