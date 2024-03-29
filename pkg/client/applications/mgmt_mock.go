// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package applications

import (
	"context"
	"github.com/auth0/go-auth0/management"
	"sync"
)

// Ensure, that applicationManagementMock does implement applicationManagement.
// If this is not the case, regenerate this file with moq.
var _ applicationManagement = &applicationManagementMock{}

// applicationManagementMock is a mock implementation of applicationManagement.
//
//	func TestSomethingThatUsesapplicationManagement(t *testing.T) {
//
//		// make and configure a mocked applicationManagement
//		mockedapplicationManagement := &applicationManagementMock{
//			ListFunc: func(ctx context.Context, opts ...management.RequestOption) (*management.ClientList, error) {
//				panic("mock out the List method")
//			},
//		}
//
//		// use mockedapplicationManagement in code that requires applicationManagement
//		// and then make assertions.
//
//	}
type applicationManagementMock struct {
	// ListFunc mocks the List method.
	ListFunc func(ctx context.Context, opts ...management.RequestOption) (*management.ClientList, error)

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
func (mock *applicationManagementMock) List(ctx context.Context, opts ...management.RequestOption) (*management.ClientList, error) {
	if mock.ListFunc == nil {
		panic("applicationManagementMock.ListFunc: method is nil but applicationManagement.List was just called")
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
//	len(mockedapplicationManagement.ListCalls())
func (mock *applicationManagementMock) ListCalls() []struct {
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
