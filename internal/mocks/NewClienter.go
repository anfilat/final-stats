// Code generated by mockery v2.5.1. DO NOT EDIT.

package mocks

import (
	symo "github.com/anfilat/final-stats/internal/symo"
	mock "github.com/stretchr/testify/mock"
)

// NewClienter is an autogenerated mock type for the NewClienter type
type NewClienter struct {
	mock.Mock
}

// NewClient provides a mock function with given fields: _a0
func (_m *NewClienter) NewClient(_a0 symo.ClientData) (<-chan *symo.Stats, func(), error) {
	ret := _m.Called(_a0)

	var r0 <-chan *symo.Stats
	if rf, ok := ret.Get(0).(func(symo.ClientData) <-chan *symo.Stats); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *symo.Stats)
		}
	}

	var r1 func()
	if rf, ok := ret.Get(1).(func(symo.ClientData) func()); ok {
		r1 = rf(_a0)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(symo.ClientData) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
