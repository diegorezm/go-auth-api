package auth

import (
	"context"

	repo "github.com/diegorezm/ticketing/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

// Register implements [Service].
func (s *svc) Register(ctx context.Context, params registerParams) error {
	panic("unimplemented")
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}
