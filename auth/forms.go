package auth

import (
	"net/http"

	"github.com/mholt/binding"
)

type CreateUserForm struct {
	Username string
}

func (form *CreateUserForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.Username: binding.Field{
			Form:     "username",
			Required: true,
		},
	}
}

type OAuthCallbackForm struct {
	Code  string
	State string
}

func (form *OAuthCallbackForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.Code: binding.Field{
			Form:     "code",
			Required: true,
		},
		&form.State: binding.Field{
			Form:     "state",
			Required: true,
		},
	}
}
