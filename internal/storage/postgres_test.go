package storage

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/zsandibe/medods-service/config"
)

func TestNewPostgresDB(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
		expectDBNil bool
	}{
		{
			name: "success",
			config: &config.Config{
				Postgres: config.PostgresConfig{
					User:     "postgres",
					Password: "postgres",
					Host:     "localhost",
					Port:     "5432",
					Name:     "medods",
				},
			},
			expectError: false,
			expectDBNil: false,
		},
		{
			name: "fail to connect with invalid config",
			config: &config.Config{
				Postgres: config.PostgresConfig{
					User:     "invalid_user",
					Password: "invalid_password",
					Host:     "invalid_host",
					Port:     "5435",
					Name:     "invalid_db",
				},
			},
			expectError: true,
			expectDBNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewPostgresDB(tt.config)
			if err != nil {
				t.Fatalf("expected flag of error: %v, got %v", tt.expectError, err)
			}
			if db != nil {
				defer db.Close()
				if err := db.DB.Ping(); err != nil {
					t.Fatalf("expected to ping successfully, got %v", err)
				}
			} else {
				t.Fatalf("expected flag of db nil: %v, got %v", tt.expectDBNil, db)
			}

		})
	}

}
