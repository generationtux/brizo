package auth

import (
	"net/http"

	"github.com/mholt/binding"
)

// CreateUserForm represents a new user form
type CreateUserForm struct {
	Username string
}

// FieldMap is used to bind CreateUserForm to a request
func (form *CreateUserForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.Username: binding.Field{
			Form:     "username",
			Required: true,
		},
	}
}

// OAuthCallbackForm represents a Github OAuth callback
type OAuthCallbackForm struct {
	Code  string
	State string
}

// FieldMap is used to bind OAuthCallbackForm to a request
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
