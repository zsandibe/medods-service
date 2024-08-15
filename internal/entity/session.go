package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id                 uuid.UUID `json:"id" db:"id"`
	Guid               uuid.UUID `json:"guid" db:"guid"`
	HashedRefreshToken []byte    `json:"hashed_refresh_token" db:"refresh_token"`
	Ip                 string    `json:"ip" db:"ip"`
	CreatedAt          time.Time `json:"created_time" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_time,omitempty" db:"updated_at"`
}
