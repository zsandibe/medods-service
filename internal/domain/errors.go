package domain

import "errors"

var (
	ErrCreatingSessions = errors.New("cannot create session")
	ErrSessionNotFound  = errors.New("session not found")
)
