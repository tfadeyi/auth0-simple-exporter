package applications

import "github.com/auth0/go-auth0/management"

//go:generate moq -out mgmt_mock.go . applicationManagement

type (
	applicationManagement interface {
		List(opts ...management.RequestOption) (l *management.ClientList, err error)
	}
)
