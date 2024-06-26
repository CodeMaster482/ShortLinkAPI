// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/usecase/link.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"
	time "time"

	model "github.com/CodeMaster482/ShortLinkAPI/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockLinkRepository is a mock of LinkRepository interface.
type MockLinkRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLinkRepositoryMockRecorder
}

// MockLinkRepositoryMockRecorder is the mock recorder for MockLinkRepository.
type MockLinkRepositoryMockRecorder struct {
	mock *MockLinkRepository
}

// NewMockLinkRepository creates a new mock instance.
func NewMockLinkRepository(ctrl *gomock.Controller) *MockLinkRepository {
	mock := &MockLinkRepository{ctrl: ctrl}
	mock.recorder = &MockLinkRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLinkRepository) EXPECT() *MockLinkRepositoryMockRecorder {
	return m.recorder
}

// GetLink mocks base method.
func (m *MockLinkRepository) GetLink(ctx context.Context, token string) (*model.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLink", ctx, token)
	ret0, _ := ret[0].(*model.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLink indicates an expected call of GetLink.
func (mr *MockLinkRepositoryMockRecorder) GetLink(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLink", reflect.TypeOf((*MockLinkRepository)(nil).GetLink), ctx, token)
}

// GetLinkByOriginal mocks base method.
func (m *MockLinkRepository) GetLinkByOriginal(ctx context.Context, origLink string) (*model.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinkByOriginal", ctx, origLink)
	ret0, _ := ret[0].(*model.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinkByOriginal indicates an expected call of GetLinkByOriginal.
func (mr *MockLinkRepositoryMockRecorder) GetLinkByOriginal(ctx, origLink interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinkByOriginal", reflect.TypeOf((*MockLinkRepository)(nil).GetLinkByOriginal), ctx, origLink)
}

// StartRecalculation mocks base method.
func (m *MockLinkRepository) StartRecalculation(interval time.Duration, deleted chan []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartRecalculation", interval, deleted)
}

// StartRecalculation indicates an expected call of StartRecalculation.
func (mr *MockLinkRepositoryMockRecorder) StartRecalculation(interval, deleted interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartRecalculation", reflect.TypeOf((*MockLinkRepository)(nil).StartRecalculation), interval, deleted)
}

// StoreLink mocks base method.
func (m *MockLinkRepository) StoreLink(ctx context.Context, link *model.Link) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreLink", ctx, link)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreLink indicates an expected call of StoreLink.
func (mr *MockLinkRepositoryMockRecorder) StoreLink(ctx, link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreLink", reflect.TypeOf((*MockLinkRepository)(nil).StoreLink), ctx, link)
}

// MockGenerator is a mock of Generator interface.
type MockGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockGeneratorMockRecorder
}

// MockGeneratorMockRecorder is the mock recorder for MockGenerator.
type MockGeneratorMockRecorder struct {
	mock *MockGenerator
}

// NewMockGenerator creates a new mock instance.
func NewMockGenerator(ctrl *gomock.Controller) *MockGenerator {
	mock := &MockGenerator{ctrl: ctrl}
	mock.recorder = &MockGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenerator) EXPECT() *MockGeneratorMockRecorder {
	return m.recorder
}

// GenerateShortURL mocks base method.
func (m *MockGenerator) GenerateShortURL(url string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateShortURL", url)
	ret0, _ := ret[0].(string)
	return ret0
}

// GenerateShortURL indicates an expected call of GenerateShortURL.
func (mr *MockGeneratorMockRecorder) GenerateShortURL(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateShortURL", reflect.TypeOf((*MockGenerator)(nil).GenerateShortURL), url)
}
