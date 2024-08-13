package v1

import "github.com/zsandibe/medods-service/internal/service"

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}
