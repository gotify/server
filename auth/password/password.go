package password

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func ValidateNewPassword(pw string) error {
	if pw == "" {
		return errors.New("password must not be empty")
	}
	if len([]byte(pw)) > 72 {
		return bcrypt.ErrPasswordTooLong
	}
	return nil
}

// CreatePassword returns a hashed version of the given password.
func CreatePassword(pw string, strength int) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), strength)
	return hashedPassword, err
}

// ComparePassword compares a hashed password with its possible plaintext equivalent.
func ComparePassword(hashedPassword, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}
