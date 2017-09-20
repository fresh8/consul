// Automatically generated by MockGen. DO NOT EDIT!
// Source: ./api_if/catalog.go

package mock_api

import (
	gomock "github.com/golang/mock/gomock"
	api "github.com/hashicorp/consul/api"
)

// Mock of Catalog interface
type MockCatalog struct {
	ctrl     *gomock.Controller
	recorder *_MockCatalogRecorder
}

// Recorder for MockCatalog (not exported)
type _MockCatalogRecorder struct {
	mock *MockCatalog
}

func NewMockCatalog(ctrl *gomock.Controller) *MockCatalog {
	mock := &MockCatalog{ctrl: ctrl}
	mock.recorder = &_MockCatalogRecorder{mock}
	return mock
}

func (_m *MockCatalog) EXPECT() *_MockCatalogRecorder {
	return _m.recorder
}

func (_m *MockCatalog) Datacenters() ([]string, error) {
	ret := _m.ctrl.Call(_m, "Datacenters")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCatalogRecorder) Datacenters() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Datacenters")
}

func (_m *MockCatalog) Deregister(_param0 *api.CatalogDeregistration, _param1 *api.WriteOptions) (*api.WriteMeta, error) {
	ret := _m.ctrl.Call(_m, "Deregister", _param0, _param1)
	ret0, _ := ret[0].(*api.WriteMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCatalogRecorder) Deregister(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Deregister", arg0, arg1)
}

func (_m *MockCatalog) Node(_param0 string, _param1 *api.QueryOptions) (*api.CatalogNode, *api.QueryMeta, error) {
	ret := _m.ctrl.Call(_m, "Node", _param0, _param1)
	ret0, _ := ret[0].(*api.CatalogNode)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockCatalogRecorder) Node(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Node", arg0, arg1)
}

func (_m *MockCatalog) Nodes(_param0 *api.QueryOptions) ([]*api.Node, *api.QueryMeta, error) {
	ret := _m.ctrl.Call(_m, "Nodes", _param0)
	ret0, _ := ret[0].([]*api.Node)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockCatalogRecorder) Nodes(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Nodes", arg0)
}

func (_m *MockCatalog) Register(_param0 *api.CatalogRegistration, _param1 *api.WriteOptions) (*api.WriteMeta, error) {
	ret := _m.ctrl.Call(_m, "Register", _param0, _param1)
	ret0, _ := ret[0].(*api.WriteMeta)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockCatalogRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Register", arg0, arg1)
}

func (_m *MockCatalog) Service(_param0 string, _param1 string, _param2 *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error) {
	ret := _m.ctrl.Call(_m, "Service", _param0, _param1, _param2)
	ret0, _ := ret[0].([]*api.CatalogService)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockCatalogRecorder) Service(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Service", arg0, arg1, arg2)
}

func (_m *MockCatalog) Services(_param0 *api.QueryOptions) (map[string][]string, *api.QueryMeta, error) {
	ret := _m.ctrl.Call(_m, "Services", _param0)
	ret0, _ := ret[0].(map[string][]string)
	ret1, _ := ret[1].(*api.QueryMeta)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

func (_mr *_MockCatalogRecorder) Services(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Services", arg0)
}
