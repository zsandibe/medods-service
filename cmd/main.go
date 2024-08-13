package main

import (
	"github.com/zsandibe/medods-service/internal/app"
	logger "github.com/zsandibe/medods-service/pkg"
)

func main() {
	if err := app.Start(); err != nil {
		logger.Error(err)
		return
	}
}
