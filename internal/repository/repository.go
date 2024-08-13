package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zsandibe/medods-service/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Create(session entity.Session) error
	GetSessionById(sessionId uuid.UUID) (entity.Session, error)
	DeleteSessionById(sessionId uuid.UUID) error
	Update(session entity.Session) error
	GetAllSessions() ([]entity.Session, error)
}

type repositoryPostgres struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repositoryPostgres {
	return &repositoryPostgres{
		db: db,
	}
}
