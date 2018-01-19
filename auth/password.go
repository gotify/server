package auth

import "golang.org/x/crypto/bcrypt"

var strength = 13

func CreatePassword(pw string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), strength)
	if err != nil {
		panic(err)
	}
	return hashedPassword
}

func ComparePassword(hashedPassword, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}
