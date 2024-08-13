package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/zsandibe/medods-service/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, session entity.Session) error
	GetSessionById(ctx context.Context, sessionId uuid.UUID) (*entity.Session, error)
	DeleteSessionById(ctx context.Context, sessionId uuid.UUID) error
	Update(ctx context.Context, session entity.Session) error
	GetAllSessions(ctx context.Context) ([]*entity.Session, error)
}

type repositoryPostgres struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repositoryPostgres {
	return &repositoryPostgres{
		db: db,
	}
}
