package auth

import "context"

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerParams struct {
	Name string `json:"name"`
	loginParams
}

type Service interface {
	Register(ctx context.Context, params registerParams) error
}
