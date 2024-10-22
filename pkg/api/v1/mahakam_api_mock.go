// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/api/v1/types.go

// Package v1 is a generated GoMock package.
package v1

import (
	gomock "github.com/golang/mock/gomock"
	apps "github.com/mahakamcloud/mahakam/pkg/api/v1/client/apps"
	clusters "github.com/mahakamcloud/mahakam/pkg/api/v1/client/clusters"
	reflect "reflect"
)

// MockClusterAPI is a mock of ClusterAPI interface
type MockClusterAPI struct {
	ctrl     *gomock.Controller
	recorder *MockClusterAPIMockRecorder
}

// MockClusterAPIMockRecorder is the mock recorder for MockClusterAPI
type MockClusterAPIMockRecorder struct {
	mock *MockClusterAPI
}

// NewMockClusterAPI creates a new mock instance
func NewMockClusterAPI(ctrl *gomock.Controller) *MockClusterAPI {
	mock := &MockClusterAPI{ctrl: ctrl}
	mock.recorder = &MockClusterAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClusterAPI) EXPECT() *MockClusterAPIMockRecorder {
	return m.recorder
}

// CreateCluster mocks base method
func (m *MockClusterAPI) CreateCluster(arg0 *clusters.CreateClusterParams) (*clusters.CreateClusterCreated, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCluster", arg0)
	ret0, _ := ret[0].(*clusters.CreateClusterCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCluster indicates an expected call of CreateCluster
func (mr *MockClusterAPIMockRecorder) CreateCluster(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCluster", reflect.TypeOf((*MockClusterAPI)(nil).CreateCluster), arg0)
}

// DescribeClusters mocks base method
func (m *MockClusterAPI) DescribeClusters(arg0 *clusters.DescribeClustersParams) (*clusters.DescribeClustersOK, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeClusters", arg0)
	ret0, _ := ret[0].(*clusters.DescribeClustersOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeClusters indicates an expected call of DescribeClusters
func (mr *MockClusterAPIMockRecorder) DescribeClusters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeClusters", reflect.TypeOf((*MockClusterAPI)(nil).DescribeClusters), arg0)
}

// GetClusters mocks base method
func (m *MockClusterAPI) GetClusters(arg0 *clusters.GetClustersParams) (*clusters.GetClustersOK, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClusters", arg0)
	ret0, _ := ret[0].(*clusters.GetClustersOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClusters indicates an expected call of GetClusters
func (mr *MockClusterAPIMockRecorder) GetClusters(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClusters", reflect.TypeOf((*MockClusterAPI)(nil).GetClusters), arg0)
}

// ValidateCluster mocks base method
func (m *MockClusterAPI) ValidateCluster(arg0 *clusters.ValidateClusterParams) (*clusters.ValidateClusterCreated, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCluster", arg0)
	ret0, _ := ret[0].(*clusters.ValidateClusterCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ValidateCluster indicates an expected call of ValidateCluster
func (mr *MockClusterAPIMockRecorder) ValidateCluster(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCluster", reflect.TypeOf((*MockClusterAPI)(nil).ValidateCluster), arg0)
}

// MockAppAPI is a mock of AppAPI interface
type MockAppAPI struct {
	ctrl     *gomock.Controller
	recorder *MockAppAPIMockRecorder
}

// MockAppAPIMockRecorder is the mock recorder for MockAppAPI
type MockAppAPIMockRecorder struct {
	mock *MockAppAPI
}

// NewMockAppAPI creates a new mock instance
func NewMockAppAPI(ctrl *gomock.Controller) *MockAppAPI {
	mock := &MockAppAPI{ctrl: ctrl}
	mock.recorder = &MockAppAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppAPI) EXPECT() *MockAppAPIMockRecorder {
	return m.recorder
}

// CreateApp mocks base method
func (m *MockAppAPI) CreateApp(arg0 *apps.CreateAppParams) (*apps.CreateAppCreated, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateApp", arg0)
	ret0, _ := ret[0].(*apps.CreateAppCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateApp indicates an expected call of CreateApp
func (mr *MockAppAPIMockRecorder) CreateApp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateApp", reflect.TypeOf((*MockAppAPI)(nil).CreateApp), arg0)
}

// GetApps mocks base method
func (m *MockAppAPI) GetApps(arg0 *apps.GetAppsParams) (*apps.GetAppsOK, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApps", arg0)
	ret0, _ := ret[0].(*apps.GetAppsOK)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApps indicates an expected call of GetApps
func (mr *MockAppAPIMockRecorder) GetApps(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApps", reflect.TypeOf((*MockAppAPI)(nil).GetApps), arg0)
}

// UploadAppValues mocks base method
func (m *MockAppAPI) UploadAppValues(arg0 *apps.UploadAppValuesParams) (*apps.UploadAppValuesCreated, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadAppValues", arg0)
	ret0, _ := ret[0].(*apps.UploadAppValuesCreated)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadAppValues indicates an expected call of UploadAppValues
func (mr *MockAppAPIMockRecorder) UploadAppValues(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadAppValues", reflect.TypeOf((*MockAppAPI)(nil).UploadAppValues), arg0)
}
