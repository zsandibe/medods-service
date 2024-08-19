package main

import (
	"github.com/zsandibe/medods-service/internal/app"
	logger "github.com/zsandibe/medods-service/pkg"
)

// @title Medods test task
// @version 1.0
// @description This is basic server for a generating JWT tokens.
// @host 0.0.0.0:7777
// @BasePath /api/v1/auth
func main() {
	if err := app.Start(); err != nil {
		logger.Error(err)
		return
	}
}
