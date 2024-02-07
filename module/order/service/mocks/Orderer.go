// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	context "context"
	models "projects/LDmitryLD/geotask/module/order/models"

	mock "github.com/stretchr/testify/mock"
)

// Orderer is an autogenerated mock type for the Orderer type
type Orderer struct {
	mock.Mock
}

// GenerateOrder provides a mock function with given fields: ctx
func (_m *Orderer) GenerateOrder(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByRadius provides a mock function with given fields: ctx, lng, lat, radius, unit
func (_m *Orderer) GetByRadius(ctx context.Context, lng float64, lat float64, radius float64, unit string) ([]models.Order, error) {
	ret := _m.Called(ctx, lng, lat, radius, unit)

	var r0 []models.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, float64, float64, float64, string) ([]models.Order, error)); ok {
		return rf(ctx, lng, lat, radius, unit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, float64, float64, float64, string) []models.Order); ok {
		r0 = rf(ctx, lng, lat, radius, unit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, float64, float64, float64, string) error); ok {
		r1 = rf(ctx, lng, lat, radius, unit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCount provides a mock function with given fields: ctx
func (_m *Orderer) GetCount(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveOldOrders provides a mock function with given fields: ctx
func (_m *Orderer) RemoveOldOrders(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveOrder provides a mock function with given fields: ctx, order
func (_m *Orderer) RemoveOrder(ctx context.Context, order models.Order) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Save provides a mock function with given fields: ctx, order
func (_m *Orderer) Save(ctx context.Context, order models.Order) error {
	ret := _m.Called(ctx, order)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Order) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewOrderer creates a new instance of Orderer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOrderer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Orderer {
	mock := &Orderer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}