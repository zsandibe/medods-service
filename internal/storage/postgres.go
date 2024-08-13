package storage

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zsandibe/medods-service/config"
	logger "github.com/zsandibe/medods-service/pkg"
)

type database struct {
	DB *sqlx.DB
}

func NewPostgresDB(cfg *config.Config) (*database, error) {
	logger.Debug("NewPostgresDB func")
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Name,
	)
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info(connString)

	err = db.Ping()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &database{
		DB: db,
	}, nil
}

func (d *database) Close() error {
	return d.DB.Close()
}

func (d *database) MigrateUp() error {
	driver, err := postgres.WithInstance(d.DB.DB, &postgres.Config{})
	if err != nil {
		logger.Error(err)
		return err
	}
	mig, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver)
	if err != nil {
		logger.Error(err)
		return err
	}
	return mig.Up()
}
