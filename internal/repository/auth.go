package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	logger "github.com/zsandibe/medods-service/pkg"

	"github.com/zsandibe/medods-service/internal/domain"
	"github.com/zsandibe/medods-service/internal/entity"

	"github.com/google/uuid"
)

func (r *repositoryPostgres) Create(ctx context.Context, session entity.Session) error {
	query := `
		INSERT INTO sessions (guid,refresh_token,ip,created_at,updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
	`
	startTime := time.Now()
	var updatedAt sql.NullTime

	_, err := r.db.ExecContext(ctx, query, &session.Guid, &session.HashedRefreshToken, &session.Ip, &startTime, &updatedAt)
	if err != nil {
		logger.Errorf("Error in inserting message: %v", err)
		return domain.ErrCreatingSessions
	}

	return nil
}

func (r *repositoryPostgres) GetSessionById(ctx context.Context, sessionId uuid.UUID) (*entity.Session, error) {
	var session entity.Session
	var updatedAt sql.NullTime

	query := `
		SELECT s.id,s.guid,
		s.refresh_token,
		s.ip,
		s.created_at,
		s.updated_at
		FROM sessions s 
		WHERE s.id = $1
	`

	if err := r.db.QueryRowContext(ctx, query, sessionId).Scan(
		&session.Id,
		&session.Guid,
		&session.HashedRefreshToken,
		&session.Ip,
		&session.CreatedAt,
		&updatedAt,
	); err != nil {
		logger.Error(err)
		return &session, err
	}

	session.UpdatedAt = updatedAt.Time

	return &session, nil
}

func (r *repositoryPostgres) DeleteSessionById(ctx context.Context, sessionId uuid.UUID) error {
	query := `
		DELETE FROM sessions WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, sessionId)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("error with executing query: %v", err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Error(ctx, fmt.Errorf("error with getting rows affected: %v", err))
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repositoryPostgres) Update(ctx context.Context, session entity.Session) error {
	query := `
	UPDATE sessions SET
	guid = $1,
	refresh_token = $2,
    updated_at= NOW()
	WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query,
		session.Guid, session.HashedRefreshToken, session.Id)
	if err != nil {
		logger.Error(ctx, fmt.Errorf("error with executing query: %v", err))
		return err
	}

	return nil
}

func (r *repositoryPostgres) GetAllSessions(ctx context.Context) ([]*entity.Session, error) {
	query := `
	SELECT id,guid,
	refresh_token,
	ip,created_at,updated_at
	FROM sessions
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Errorf("error querying with context: %v", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []*entity.Session
	for rows.Next() {
		session := &entity.Session{}

		var updatedAt sql.NullTime
		if err := rows.Scan(
			&session.Id,
			&session.Guid,
			&session.HashedRefreshToken,
			&session.Ip,
			&session.CreatedAt,
			&updatedAt,
		); err != nil {
			logger.Errorf("problems with scanning rows: %v", err)
			return sessions, err
		}
		session.UpdatedAt = updatedAt.Time

		sessions = append(sessions, session)
	}

	if err := rows.Close(); err != nil {
		logger.Error(err)
		return sessions, err
	}

	if err := rows.Err(); err != nil {
		logger.Error(err)
		return sessions, err
	}

	logger.Debugf("sessions: %+v", sessions)
	return sessions, nil
}
