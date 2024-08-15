package service

import (
	"context"

	"github.com/zsandibe/medods-service/config"
	"github.com/zsandibe/medods-service/internal/domain"
	"github.com/zsandibe/medods-service/internal/entity"
	"github.com/zsandibe/medods-service/internal/repository"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, guid uuid.UUID, ip string) (domain.TokenPair, error)
	Update(ctx context.Context, sessionId uuid.UUID) (domain.TokenPair, error)
	GetAllSessions(ctx context.Context) ([]*entity.Session, error)
	NotifyToEmail(oldIp, newIp string) error
	GetSessionById(ctx context.Context, sessionId uuid.UUID) (*entity.Session, error)
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
