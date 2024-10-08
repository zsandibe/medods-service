package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zsandibe/medods-service/pkg"
	logger "github.com/zsandibe/medods-service/pkg"

	"github.com/zsandibe/medods-service/internal/domain"
	"github.com/zsandibe/medods-service/internal/entity"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func (s *service) GetAllSessions(ctx context.Context) ([]*entity.Session, error) {
	logger.Info("service.GetAllSessions(): getting all sessions")
	return s.repo.GetAllSessions(ctx)
}

func (s *service) Create(ctx context.Context, guid uuid.UUID, ip string) (domain.TokenPair, error) {
	logger.Debugf("service.Create():\nBefore: guid %v", guid)
	session := entity.Session{
		Id:        uuid.New(),
		Guid:      guid,
		Ip:        ip,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tokenPair, err := s.signTokenPair(session.Id, session.Guid)
	if err != nil {
		return domain.TokenPair{}, err
	}

	session.HashedRefreshToken, err = pkg.GetHashFromToken(tokenPair.RefreshToken)
	if err != nil {
		return domain.TokenPair{}, err
	}

	if err = s.repo.Create(ctx, session); err != nil {
		return domain.TokenPair{}, err
	}
	logger.Debugf("service.Create():\nAfter: token pair %v", tokenPair)

	logger.Info("service.Create(): token pairs created")
	return tokenPair, nil
}

func (s *service) Update(ctx context.Context, sessionId uuid.UUID) (domain.TokenPair, error) {
	logger.Debugf("service.Update():\nBefore: session id %v", sessionId)
	session, err := s.repo.GetSessionById(ctx, sessionId)
	if err != nil {
		logger.Errorf("Error getting session: %v", err)
		return domain.TokenPair{}, domain.ErrSessionNotFound
	}
	fmt.Println(session.HashedRefreshToken, "before")
	if session.UpdatedAt.Sub(time.Now()) >= s.conf.Token.RefreshTokenAge {
		if err = s.repo.DeleteSessionById(ctx, sessionId); err != nil {
			return domain.TokenPair{}, err
		}
		return domain.TokenPair{}, errors.New("refresh token expired")
	}

	tokenPair, err := s.signTokenPair(session.Id, session.Guid)
	if err != nil {
		return domain.TokenPair{}, err
	}

	session.UpdatedAt = time.Now()

	session.HashedRefreshToken, err = pkg.GetHashFromToken(tokenPair.RefreshToken)
	if err != nil {
		return domain.TokenPair{}, err
	}

	if err = s.repo.Update(ctx, *session); err != nil {
		return domain.TokenPair{}, err
	}

	logger.Debugf("service.Update():\nAfter: updated token pairs %v", tokenPair)
	logger.Info("service.Update(): token pairs successfully updated")
	return tokenPair, nil
}

func (s *service) GetSessionById(ctx context.Context, sessionId uuid.UUID) (*entity.Session, error) {
	session, err := s.repo.GetSessionById(ctx, sessionId)
	if err != nil {
		return nil, domain.ErrSessionNotFound
	}
	return session, nil
}

func (s *service) signTokenPair(sessionId uuid.UUID, guid uuid.UUID) (domain.TokenPair, error) {
	logger.Debugf("service.signTokenPair():\nBefore: sessionId - %v, guid - %v", sessionId, guid)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, domain.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.conf.Token.AccessTokenAge)),
		},
		SessionId: sessionId,
		Guid:      guid,
	})

	signedAccessToken, err := accessToken.SignedString([]byte(s.conf.Token.AccessKey))
	if err != nil {
		fmt.Println(err)
		return domain.TokenPair{}, errors.New("failed to signed")
	}
	refreshToken := uuid.New()

	tokenPair := domain.TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: refreshToken[:],
	}
	logger.Debugf("service.signTokenPair():\nAfter: token pairs - %v", err)
	logger.Info("service.ignTokenPair(): token pairs successfully signed")
	return tokenPair, nil
}
