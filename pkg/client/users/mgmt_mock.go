// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package users

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"sync"
)

// Ensure, that userManagementMock does implement userManagement.
// If this is not the case, regenerate this file with moq.
var _ userManagement = &userManagementMock{}

// userManagementMock is a mock implementation of userManagement.
//
//	func TestSomethingThatUsesuserManagement(t *testing.T) {
//
//		// make and configure a mocked userManagement
//		mockeduserManagement := &userManagementMock{
//			ListFunc: func(ctx context.Context, opts ...management.RequestOption) (*management.UserList, error) {
//				panic("mock out the List method")
//			},
//		}
//
//		// use mockeduserManagement in code that requires userManagement
//		// and then make assertions.
//
//	}
type userManagementMock struct {
	// ListFunc mocks the List method.
	ListFunc func(ctx context.Context, opts ...management.RequestOption) (*management.UserList, error)

	// calls tracks calls to the methods.
	calls struct {
		// List holds details about calls to the List method.
		List []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Opts is the opts argument value.
			Opts []management.RequestOption
		}
	}
	lockList sync.RWMutex
}

// List calls ListFunc.
func (mock *userManagementMock) List(ctx context.Context, opts ...management.RequestOption) (*management.UserList, error) {
	if mock.ListFunc == nil {
		panic("userManagementMock.ListFunc: method is nil but userManagement.List was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Opts []management.RequestOption
	}{
		Ctx:  ctx,
		Opts: opts,
	}
	mock.lockList.Lock()
	mock.calls.List = append(mock.calls.List, callInfo)
	mock.lockList.Unlock()
	return mock.ListFunc(ctx, opts...)
}

// ListCalls gets all the calls that were made to List.
// Check the length with:
//
//	len(mockeduserManagement.ListCalls())
func (mock *userManagementMock) ListCalls() []struct {
	Ctx  context.Context
	Opts []management.RequestOption
} {
	var calls []struct {
		Ctx  context.Context
		Opts []management.RequestOption
	}
	mock.lockList.RLock()
	calls = mock.calls.List
	mock.lockList.RUnlock()
	return calls
}