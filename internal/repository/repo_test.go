package repository

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zsandibe/medods-service/internal/entity"
)

func Test_repositoryPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewRepository(sqlxDB)

	tests := []struct {
		name    string
		session entity.Session
		wantErr error
		mock    func(session entity.Session)
	}{
		{
			name: "success",
			session: entity.Session{
				Guid:               uuid.New(),
				HashedRefreshToken: []byte("388d6b48-1a7e-4df7-b1eb-79ca726fb814"),
				Ip:                 "127.0.0.1",
			},
			wantErr: nil,
			mock: func(session entity.Session) {
				mock.ExpectExec(`INSERT INTO sessions`).
					WithArgs(session.Guid, session.HashedRefreshToken, session.Ip, sqlmock.AnyArg(), sql.NullTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "failure",
			session: entity.Session{
				Guid:               uuid.New(),
				HashedRefreshToken: []byte("188d6348-1a7e-4df7-b1eb-79ca726fb814"),
				Ip:                 "127.5.0.1",
			},
			wantErr: nil,
			mock: func(session entity.Session) {
				mock.ExpectExec(`INSERT INTO sessions`).
					WithArgs(session.Guid, session.HashedRefreshToken, session.Ip, sqlmock.AnyArg(), sql.NullTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.session)
			err := repo.Create(context.Background(), tt.session)
			if err != tt.wantErr {
				t.Errorf("repositoryPostgres.Create() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repositoryPostgres_GetSessionById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewRepository(sqlxDB)

	sessionID := uuid.New()
	guid := uuid.New()
	hashedToken := []byte("188d6348-1a7e-4df7-b1eb-79ca726fb814")
	ip := "127.0.0.1"
	createdAt := time.Now()
	updatedAt := time.Now()

	tests := []struct {
		name      string
		sessionId uuid.UUID
		want      *entity.Session
		wantErr   error
		mock      func(sessionId uuid.UUID)
	}{
		{
			name:      "success",
			sessionId: sessionID,
			want: &entity.Session{
				Id:                 sessionID,
				Guid:               guid,
				HashedRefreshToken: hashedToken,
				Ip:                 ip,
				CreatedAt:          createdAt,
				UpdatedAt:          updatedAt,
			},
			wantErr: nil,
			mock: func(sessionId uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "guid", "refresh_token", "ip", "created_at", "updated_at"}).
					AddRow(sessionID, guid, hashedToken, ip, createdAt, updatedAt)
				mock.ExpectQuery(`SELECT s\.id,\s*s\.guid,\s*s\.refresh_token,\s*s\.ip,\s*s\.created_at,\s*s\.updated_at FROM sessions s WHERE s\.id = \$1`).
					WithArgs(sessionId).
					WillReturnRows(rows)
			},
		},
		{
			name:      "not found",
			sessionId: sessionID,
			want:      nil,
			wantErr:   sql.ErrNoRows,
			mock: func(sessionId uuid.UUID) {
				rows := sqlmock.NewRows([]string{"id", "guid", "refresh_token", "ip", "created_at", "updated_at"}).
					AddRow(uuid.New(), guid, hashedToken, ip, createdAt, updatedAt)
				mock.ExpectQuery(`SELECT s\.id,\s*s\.guid,\s*s\.refresh_token,\s*s\.ip,\s*s\.created_at,\s*s\.updated_at FROM sessions s WHERE s\.id = \$1`).
					WithArgs(sessionId).
					WillReturnRows(rows).WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.sessionId)

			got, err := repo.GetSessionById(context.Background(), tt.sessionId)
			if err != tt.wantErr {
				t.Errorf("repositoryPostgres.GetSessionById() error = %v, want %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != nil && tt.want != nil {
				if !got.CreatedAt.Equal(tt.want.CreatedAt) || !got.UpdatedAt.Equal(tt.want.UpdatedAt) {
					t.Errorf("repositoryPostgres.GetSessionById() times do not match, got = %v, want %v", got, tt.want)
				}

				got.CreatedAt, got.UpdatedAt = tt.want.CreatedAt, tt.want.UpdatedAt
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("repositoryPostgres.GetSessionById() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_repositoryPostgres_DeleteSessionById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewRepository(sqlxDB)

	sessionID := uuid.New()

	tests := []struct {
		name      string
		sessionId uuid.UUID
		wantErr   error
		mock      func(sessionId uuid.UUID)
	}{
		{
			name:      "success",
			sessionId: sessionID,
			wantErr:   nil,
			mock: func(sessionId uuid.UUID) {
				mock.ExpectExec(`DELETE FROM sessions WHERE id = \$1`).
					WithArgs(sessionId).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:      "not found",
			sessionId: sessionID,
			wantErr:   sql.ErrNoRows,
			mock: func(sessionId uuid.UUID) {
				mock.ExpectExec(`DELETE FROM sessions WHERE id = \$1`).
					WithArgs(sessionId).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.sessionId)

			err := repo.DeleteSessionById(context.Background(), tt.sessionId)
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("repositoryPostgres.DeleteSessionById() error = %v, want %v", err, tt.wantErr)
			} else if err == nil && tt.wantErr != nil {
				t.Errorf("repositoryPostgres.DeleteSessionById() expected error = %v, got nil", tt.wantErr)
			}
		})
	}
}

func Test_repositoryPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock db: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewRepository(sqlxDB)

	sessionID := uuid.New()
	guid := uuid.New()
	hashedToken := []byte("188d6348-1a7e-4df7-b1eb-79ca726fb814")

	tests := []struct {
		name    string
		session entity.Session
		wantErr error
		mock    func(session entity.Session)
	}{
		{
			name: "success",
			session: entity.Session{
				Id:                 sessionID,
				Guid:               guid,
				HashedRefreshToken: hashedToken,
			},
			mock: func(session entity.Session) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE sessions SET guid = $1, refresh_token = $2, updated_at= NOW() WHERE id = $3")).
					WithArgs(session.Guid, session.HashedRefreshToken, session.Id).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: nil,
		},
		{
			name: "not found",
			session: entity.Session{
				Id:                 sessionID,
				Guid:               guid,
				HashedRefreshToken: hashedToken,
			},
			mock: func(session entity.Session) {
				mock.ExpectExec(regexp.QuoteMeta("UPDATE sessions SET guid = $1, refresh_token = $2, updated_at= NOW() WHERE id = $3")).
					WithArgs(session.Guid, session.HashedRefreshToken, session.Id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.session)

			err := repo.Update(context.Background(), tt.session)
			if err == nil && tt.wantErr == sql.ErrNoRows {
				err = sql.ErrNoRows
			}

			if err != tt.wantErr {
				t.Errorf("repositoryPostgres.Update() error = %v, want %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet mock expectations: %v", err)
			}
		})
	}
}
