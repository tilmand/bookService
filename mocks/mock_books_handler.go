// Code generated by MockGen. DO NOT EDIT.
// Source: bookService/http (interfaces: BooksHandlerInterface)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockBooksHandlerInterface is a mock of BooksHandlerInterface interface.
type MockBooksHandlerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockBooksHandlerInterfaceMockRecorder
}

// MockBooksHandlerInterfaceMockRecorder is the mock recorder for MockBooksHandlerInterface.
type MockBooksHandlerInterfaceMockRecorder struct {
	mock *MockBooksHandlerInterface
}

// NewMockBooksHandlerInterface creates a new mock instance.
func NewMockBooksHandlerInterface(ctrl *gomock.Controller) *MockBooksHandlerInterface {
	mock := &MockBooksHandlerInterface{ctrl: ctrl}
	mock.recorder = &MockBooksHandlerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBooksHandlerInterface) EXPECT() *MockBooksHandlerInterfaceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockBooksHandlerInterface) Add(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", arg0)
}

// Add indicates an expected call of Add.
func (mr *MockBooksHandlerInterfaceMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockBooksHandlerInterface)(nil).Add), arg0)
}

// Delete mocks base method.
func (m *MockBooksHandlerInterface) Delete(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", arg0)
}

// Delete indicates an expected call of Delete.
func (mr *MockBooksHandlerInterfaceMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBooksHandlerInterface)(nil).Delete), arg0)
}

// Find mocks base method.
func (m *MockBooksHandlerInterface) Find(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Find", arg0)
}

// Find indicates an expected call of Find.
func (mr *MockBooksHandlerInterfaceMockRecorder) Find(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockBooksHandlerInterface)(nil).Find), arg0)
}

// GetAll mocks base method.
func (m *MockBooksHandlerInterface) GetAll(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GetAll", arg0)
}

// GetAll indicates an expected call of GetAll.
func (mr *MockBooksHandlerInterfaceMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockBooksHandlerInterface)(nil).GetAll), arg0)
}

// Update mocks base method.
func (m *MockBooksHandlerInterface) Update(arg0 *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0)
}

// Update indicates an expected call of Update.
func (mr *MockBooksHandlerInterfaceMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockBooksHandlerInterface)(nil).Update), arg0)
}
