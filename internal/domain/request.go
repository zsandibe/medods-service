package domain

import "github.com/google/uuid"

type LoginRequest struct {
	Guid uuid.UUID `json:"guid"`
}

type RefreshRequest struct {
	SessionID uuid.UUID `json:"session_id"`
}
