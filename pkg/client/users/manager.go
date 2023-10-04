package users

import (
	"context"
	"github.com/auth0/go-auth0/management"
)

//go:generate moq -out mgmt_mock.go . userManagement

type (
	userManagement interface {
		List(ctx context.Context, opts ...management.RequestOption) (l *management.UserList, err error)
	}
)
