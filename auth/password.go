package auth

import "golang.org/x/crypto/bcrypt"

var strength = 13

// CreatePassword returns a hashed version of the given password.
func CreatePassword(pw string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), strength)
	if err != nil {
		panic(err)
	}
	return hashedPassword
}

// ComparePassword compares a hashed password with its possible plaintext equivalent.
func ComparePassword(hashedPassword, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}
