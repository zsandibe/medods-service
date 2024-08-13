package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id                 uuid.UUID `bson:"id" json:"id"`
	Guid               uuid.UUID `bson:"guid" json:"guid"`
	HashedRefreshToken []byte    `bson:"hashed_refresh_token" json:"hashed_refresh_token"`
	CreatedAt          time.Time `bson:"created_at" json:"created_time"`
	UpdatedAt          time.Time `bson:"updated_at" json:"updated_time"`
}
