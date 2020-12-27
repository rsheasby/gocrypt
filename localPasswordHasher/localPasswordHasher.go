package localPasswordHasher

import (
	"crypto/sha512"
	"fmt"

	"github.com/rsheasby/gocrypt"
	"golang.org/x/crypto/bcrypt"
)

type localPasswordHasher struct {
	cost int
}

func New(cost int) (lph gocrypt.PasswordHasher, err error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return gocrypt.PasswordHasher(nil), fmt.Errorf(`cost %d is invalid - cost must be between %d and %d`, cost, bcrypt.MinCost, bcrypt.MaxCost)
	}
	return &localPasswordHasher{cost: cost}, nil
}

func (l *localPasswordHasher) HashPassword(password string) (hash string, err error) {
	if len(password) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}

	shaBytes := sha512.Sum512([]byte(password))
	hashBytes, err := bcrypt.GenerateFromPassword(shaBytes[:], l.cost)
	if err != nil {
		// This should never fail, except for maybe OOM errors. May as well return the error just in-case anyway though.
		return "", err
	}
	return string(hashBytes), nil
}

func (l *localPasswordHasher) ValidatePassword(password string, hash string) (isValid bool, err error) {
	if len(password) == 0 {
		return false, fmt.Errorf("password cannot be empty")
	}

	pwdHash := sha512.Sum512([]byte(password))
	err = bcrypt.CompareHashAndPassword([]byte(hash), pwdHash[:])
	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else {
		return false, err
	}
}
