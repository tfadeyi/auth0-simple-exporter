package logs

import "github.com/auth0/go-auth0/management"

//go:generate rm -f ./mgmt_mock.go
//go:generate moq -out mgmt_mock.go . logManagement

type (
	logManagement interface {
		Read(id string, opts ...management.RequestOption) (l *management.Log, err error)
		List(opts ...management.RequestOption) (l []*management.Log, err error)
		Search(opts ...management.RequestOption) ([]*management.Log, error)
	}
)
