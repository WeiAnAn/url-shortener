// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/short_url/service.go

// Package mock_shorturl is a generated GoMock package.
package mock_shorturl

import (
	context "context"
	reflect "reflect"
	time "time"

	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	gomock "github.com/golang/mock/gomock"
)

// MockShortURLGenerator is a mock of ShortURLGenerator interface.
type MockShortURLGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockShortURLGeneratorMockRecorder
}

// MockShortURLGeneratorMockRecorder is the mock recorder for MockShortURLGenerator.
type MockShortURLGeneratorMockRecorder struct {
	mock *MockShortURLGenerator
}

// NewMockShortURLGenerator creates a new mock instance.
func NewMockShortURLGenerator(ctrl *gomock.Controller) *MockShortURLGenerator {
	mock := &MockShortURLGenerator{ctrl: ctrl}
	mock.recorder = &MockShortURLGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortURLGenerator) EXPECT() *MockShortURLGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockShortURLGenerator) Generate(arg0 int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockShortURLGeneratorMockRecorder) Generate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockShortURLGenerator)(nil).Generate), arg0)
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateShortURL mocks base method.
func (m *MockService) CreateShortURL(arg0 context.Context, arg1 string, arg2 time.Time) (*shorturl.ShortURLWithExpireTime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortURL", arg0, arg1, arg2)
	ret0, _ := ret[0].(*shorturl.ShortURLWithExpireTime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortURL indicates an expected call of CreateShortURL.
func (mr *MockServiceMockRecorder) CreateShortURL(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortURL", reflect.TypeOf((*MockService)(nil).CreateShortURL), arg0, arg1, arg2)
}

// GetOriginalURL mocks base method.
func (m *MockService) GetOriginalURL(arg0 context.Context, arg1 string) (*shorturl.ShortURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURL", arg0, arg1)
	ret0, _ := ret[0].(*shorturl.ShortURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURL indicates an expected call of GetOriginalURL.
func (mr *MockServiceMockRecorder) GetOriginalURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURL", reflect.TypeOf((*MockService)(nil).GetOriginalURL), arg0, arg1)
}
