package service

import (
	"github.com/zsandibe/medods-service/config"
	"github.com/zsandibe/medods-service/internal/domain"
	"github.com/zsandibe/medods-service/internal/entity"
	"github.com/zsandibe/medods-service/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	Create(guid uuid.UUID) (domain.TokenPair, error)
	Update(sessionId uuid.UUID) (domain.TokenPair, error)
	GetAllSessions() ([]entity.Session, error)
}

type service struct {
	conf *config.Config
	repo repository.Repository
}

func NewService(repo repository.Repository, conf *config.Config) *service {
	return &service{
		repo: repo,
		conf: conf,
	}
}
