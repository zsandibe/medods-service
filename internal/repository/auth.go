package repository

import (
	"fmt"

	logger "github.com/zsandibe/medods-service/pkg"

	"github.com/zsandibe/medods-service/internal/entity"

	"github.com/google/uuid"
)

func (r *repositoryPostgres) Create(session entity.Session) error {
	return nil
}

func (r *repositoryPostgres) GetSessionById(sessionId uuid.UUID) (entity.Session, error) {
	var session entity.Session

	return session, nil
}

func (r *repositoryPostgres) DeleteSessionById(sessionId uuid.UUID) error {
	return nil
}

func (r *repositoryPostgres) Update(session entity.Session) error {
	fmt.Println(session.HashedRefreshToken, "after")

	return nil
}

func (r *repositoryPostgres) GetAllSessions() ([]entity.Session, error) {
	var sessions []entity.Session

	logger.Debugf("sessions: %v", sessions)
	return sessions, nil
}
