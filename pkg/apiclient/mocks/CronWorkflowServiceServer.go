// Code generated by mockery v1.1.1. DO NOT EDIT.

package mocks

import (
	context "context"

	cronworkflow "github.com/argoproj/argo/pkg/apiclient/cronworkflow"
	mock "github.com/stretchr/testify/mock"

	v1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
)

// CronWorkflowServiceServer is an autogenerated mock type for the CronWorkflowServiceServer type
type CronWorkflowServiceServer struct {
	mock.Mock
}

// CreateCronWorkflow provides a mock function with given fields: _a0, _a1
func (_m *CronWorkflowServiceServer) CreateCronWorkflow(_a0 context.Context, _a1 *cronworkflow.CreateCronWorkflowRequest) (*v1alpha1.CronWorkflow, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *v1alpha1.CronWorkflow
	if rf, ok := ret.Get(0).(func(context.Context, *cronworkflow.CreateCronWorkflowRequest) *v1alpha1.CronWorkflow); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.CronWorkflow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *cronworkflow.CreateCronWorkflowRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCronWorkflow provides a mock function with given fields: _a0, _a1
func (_m *CronWorkflowServiceServer) DeleteCronWorkflow(_a0 context.Context, _a1 *cronworkflow.DeleteCronWorkflowRequest) (*cronworkflow.CronWorkflowDeletedResponse, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *cronworkflow.CronWorkflowDeletedResponse
	if rf, ok := ret.Get(0).(func(context.Context, *cronworkflow.DeleteCronWorkflowRequest) *cronworkflow.CronWorkflowDeletedResponse); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cronworkflow.CronWorkflowDeletedResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *cronworkflow.DeleteCronWorkflowRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCronWorkflow provides a mock function with given fields: _a0, _a1
func (_m *CronWorkflowServiceServer) GetCronWorkflow(_a0 context.Context, _a1 *cronworkflow.GetCronWorkflowRequest) (*v1alpha1.CronWorkflow, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *v1alpha1.CronWorkflow
	if rf, ok := ret.Get(0).(func(context.Context, *cronworkflow.GetCronWorkflowRequest) *v1alpha1.CronWorkflow); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.CronWorkflow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *cronworkflow.GetCronWorkflowRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LintCronWorkflow provides a mock function with given fields: _a0, _a1
func (_m *CronWorkflowServiceServer) LintCronWorkflow(_a0 context.Context, _a1 *cronworkflow.LintCronWorkflowRequest) (*v1alpha1.CronWorkflow, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *v1alpha1.CronWorkflow
	if rf, ok := ret.Get(0).(func(context.Context, *cronworkflow.LintCronWorkflowRequest) *v1alpha1.CronWorkflow); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.CronWorkflow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *cronworkflow.LintCronWorkflowRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListCronWorkflows provides a mock function with given fields: _a0, _a1
func (_m *CronWorkflowServiceServer) ListCronWorkflows(_a0 context.Context, _a1 *cronworkflow.ListCronWorkflowsRequest) (*v1alpha1.CronWorkflowList, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *v1alpha1.CronWorkflowList
	if rf, ok := ret.Get(0).(func(context.Context, *cronworkflow.ListCronWorkflowsRequest) *v1alpha1.CronWorkflowList); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.CronWorkflowList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *cronworkflow.ListCronWorkflowsRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCronWorkflow provides a mock function with given fields: _a0, _a1
func (_m *CronWorkflowServiceServer) UpdateCronWorkflow(_a0 context.Context, _a1 *cronworkflow.UpdateCronWorkflowRequest) (*v1alpha1.CronWorkflow, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *v1alpha1.CronWorkflow
	if rf, ok := ret.Get(0).(func(context.Context, *cronworkflow.UpdateCronWorkflowRequest) *v1alpha1.CronWorkflow); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.CronWorkflow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *cronworkflow.UpdateCronWorkflowRequest) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
