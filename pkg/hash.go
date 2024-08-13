package pkg

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GetHashFromToken(token []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to generate hash from token")
	}
	return bytes, nil
}
