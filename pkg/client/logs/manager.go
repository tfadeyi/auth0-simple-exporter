package logs

import (
	"context"
	"github.com/auth0/go-auth0/management"
)

//go:generate moq -out mgmt_mock.go . logManagement

type (
	logManagement interface {
		Read(ctx context.Context, id string, opts ...management.RequestOption) (l *management.Log, err error)
		List(ctx context.Context, opts ...management.RequestOption) (l []*management.Log, err error)
		Search(ctx context.Context, opts ...management.RequestOption) ([]*management.Log, error)
	}
)
