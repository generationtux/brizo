package auth

import (
	"net/http"

	"github.com/mholt/binding"
)

// CreateUserForm represents a request to create a new user
type CreateUserForm struct {
	Username string
}

// FieldMap to map JSON request to struct
func (form *CreateUserForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.Username: binding.Field{
			Form:     "username",
			Required: true,
		},
	}
}

// OAuthCallbackForm represents OAuth response data
type OAuthCallbackForm struct {
	Code  string
	State string
}

// FieldMap to map form values to a struct
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

// ValidateJWTForm represents a json request to validate a token
type ValidateJWTForm struct {
	Token string
}

// FieldMap to map form values to a struct
func (form *ValidateJWTForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.Token: binding.Field{Form: "token", Required: true},
	}
}
