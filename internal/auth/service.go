package auth

import (
	"context"
	"errors"
	"log/slog"

	repo "github.com/diegorezm/ticketing/internal/adapters/postgresql/sqlc"
	"github.com/diegorezm/ticketing/internal/jwt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInternal           = errors.New("internal server error")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type svc struct {
	repo   *repo.Queries
	db     *pgx.Conn
	logger *slog.Logger
}

func NewService(repo *repo.Queries, db *pgx.Conn, logger *slog.Logger) Service {
	return &svc{
		repo:   repo,
		db:     db,
		logger: logger,
	}
}

func (s *svc) Register(ctx context.Context, params registerParams) (err error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("begin tx failed", "err", err)
		return ErrInternal
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	qtx := s.repo.WithTx(tx)

	passwordHash, err := HashPassword(params.Password)
	if err != nil {
		s.logger.Error("hash password failed", "err", err)
		return ErrInternal
	}

	_, err = qtx.CreateUser(ctx, repo.CreateUserParams{
		Email:        params.Email,
		PasswordHash: passwordHash,
		Name:         params.Name,
	})
	if err != nil {
		s.logger.Error("create user failed", "err", err)

		if isUniqueViolation(err) {
			return ErrEmailAlreadyExists
		}

		return ErrInternal
	}

	if err = tx.Commit(ctx); err != nil {
		s.logger.Error("commit failed", "err", err)
		return ErrInternal
	}

	return nil
}

func (s *svc) Login(ctx context.Context, params loginParams) (loginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, params.Email)

	if err != nil {
		s.logger.Error("get user failed", "err", err)
		return loginResponse{}, ErrInvalidCredentials
	}

	if !VerifyPassword(user.PasswordHash, params.Password) {
		return loginResponse{}, ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID.String())
	if err != nil {
		s.logger.Error("generate token failed", "err", err)
		return loginResponse{}, ErrInternal
	}

	return loginResponse{
		Token: token,
		User: User{
			Name:  user.Name,
			Email: user.Email,
			Id:    user.ID.String(),
		},
	}, nil
}

func (s *svc) Me(ctx context.Context, userId string) (User, error) {
	var uuid pgtype.UUID

	if err := uuid.Scan(userId); err != nil {
		return User{}, ErrInvalidCredentials
	}

	user, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		return User{}, ErrInvalidCredentials
	}

	return User{
		Name:  user.Name,
		Email: user.Email,
		Id:    user.ID.String(),
	}, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
