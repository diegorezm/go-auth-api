package auth

import "context"

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type loginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type registerParams struct {
	Name string `json:"name"`
	loginParams
}

type Service interface {
	Register(ctx context.Context, params registerParams) error
	Login(ctx context.Context, params loginParams) (loginResponse, error)
	Me(ctx context.Context, userId string) (User, error)
}
