// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/cluster/gitrepo/gitrepo_cluster.go

// Package mock_gitrepo is a generated GoMock package.
package mock_gitrepo

import (
	context "context"
	gitrepo "g.hz.netease.com/horizon/pkg/cluster/gitrepo"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClusterGitRepo is a mock of ClusterGitRepo interface
type MockClusterGitRepo struct {
	ctrl     *gomock.Controller
	recorder *MockClusterGitRepoMockRecorder
}

// MockClusterGitRepoMockRecorder is the mock recorder for MockClusterGitRepo
type MockClusterGitRepoMockRecorder struct {
	mock *MockClusterGitRepo
}

// NewMockClusterGitRepo creates a new mock instance
func NewMockClusterGitRepo(ctrl *gomock.Controller) *MockClusterGitRepo {
	mock := &MockClusterGitRepo{ctrl: ctrl}
	mock.recorder = &MockClusterGitRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClusterGitRepo) EXPECT() *MockClusterGitRepoMockRecorder {
	return m.recorder
}

// GetCluster mocks base method
func (m *MockClusterGitRepo) GetCluster(ctx context.Context, application, cluster, templateName string) (*gitrepo.ClusterFiles, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCluster", ctx, application, cluster, templateName)
	ret0, _ := ret[0].(*gitrepo.ClusterFiles)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCluster indicates an expected call of GetCluster
func (mr *MockClusterGitRepoMockRecorder) GetCluster(ctx, application, cluster, templateName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCluster", reflect.TypeOf((*MockClusterGitRepo)(nil).GetCluster), ctx, application, cluster, templateName)
}

// CreateCluster mocks base method
func (m *MockClusterGitRepo) CreateCluster(ctx context.Context, params *gitrepo.CreateClusterParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCluster", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCluster indicates an expected call of CreateCluster
func (mr *MockClusterGitRepoMockRecorder) CreateCluster(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCluster", reflect.TypeOf((*MockClusterGitRepo)(nil).CreateCluster), ctx, params)
}

// UpdateCluster mocks base method
func (m *MockClusterGitRepo) UpdateCluster(ctx context.Context, params *gitrepo.UpdateClusterParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCluster", ctx, params)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCluster indicates an expected call of UpdateCluster
func (mr *MockClusterGitRepoMockRecorder) UpdateCluster(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCluster", reflect.TypeOf((*MockClusterGitRepo)(nil).UpdateCluster), ctx, params)
}

// DeleteCluster mocks base method
func (m *MockClusterGitRepo) DeleteCluster(ctx context.Context, application, cluster string, clusterID uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCluster", ctx, application, cluster, clusterID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCluster indicates an expected call of DeleteCluster
func (mr *MockClusterGitRepoMockRecorder) DeleteCluster(ctx, application, cluster, clusterID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCluster", reflect.TypeOf((*MockClusterGitRepo)(nil).DeleteCluster), ctx, application, cluster, clusterID)
}

// CompareConfig mocks base method
func (m *MockClusterGitRepo) CompareConfig(ctx context.Context, application, cluster string, from, to *string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CompareConfig", ctx, application, cluster, from, to)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CompareConfig indicates an expected call of CompareConfig
func (mr *MockClusterGitRepoMockRecorder) CompareConfig(ctx, application, cluster, from, to interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CompareConfig", reflect.TypeOf((*MockClusterGitRepo)(nil).CompareConfig), ctx, application, cluster, from, to)
}

// MergeBranch mocks base method
func (m *MockClusterGitRepo) MergeBranch(ctx context.Context, application, cluster string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MergeBranch", ctx, application, cluster)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MergeBranch indicates an expected call of MergeBranch
func (mr *MockClusterGitRepoMockRecorder) MergeBranch(ctx, application, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MergeBranch", reflect.TypeOf((*MockClusterGitRepo)(nil).MergeBranch), ctx, application, cluster)
}

// UpdateImage mocks base method
func (m *MockClusterGitRepo) UpdateImage(ctx context.Context, application, cluster, template, image string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateImage", ctx, application, cluster, template, image)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateImage indicates an expected call of UpdateImage
func (mr *MockClusterGitRepoMockRecorder) UpdateImage(ctx, application, cluster, template, image interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateImage", reflect.TypeOf((*MockClusterGitRepo)(nil).UpdateImage), ctx, application, cluster, template, image)
}

// GetConfigCommit mocks base method
func (m *MockClusterGitRepo) GetConfigCommit(ctx context.Context, application, cluster string) (*gitrepo.ClusterCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfigCommit", ctx, application, cluster)
	ret0, _ := ret[0].(*gitrepo.ClusterCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfigCommit indicates an expected call of GetConfigCommit
func (mr *MockClusterGitRepoMockRecorder) GetConfigCommit(ctx, application, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigCommit", reflect.TypeOf((*MockClusterGitRepo)(nil).GetConfigCommit), ctx, application, cluster)
}

// GetRepoInfo mocks base method
func (m *MockClusterGitRepo) GetRepoInfo(ctx context.Context, application, cluster string) *gitrepo.RepoInfo {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepoInfo", ctx, application, cluster)
	ret0, _ := ret[0].(*gitrepo.RepoInfo)
	return ret0
}

// GetRepoInfo indicates an expected call of GetRepoInfo
func (mr *MockClusterGitRepoMockRecorder) GetRepoInfo(ctx, application, cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoInfo", reflect.TypeOf((*MockClusterGitRepo)(nil).GetRepoInfo), ctx, application, cluster)
}
