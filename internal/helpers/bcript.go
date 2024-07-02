package helpers

import "golang.org/x/crypto/bcrypt"

// to hash the password inserted by user when register: returning the generated/hashed password (string)
func HashPass(p string) string {

	cost := 8 // parameter for how many hashing is iterated
	password := []byte(p) //converting from string to slice of byte
	hash, _ := bcrypt.GenerateFromPassword(password, cost)

	return string(hash)
}

// to compare the existing password with the one inputted from user when login, returning boolean: true/false
func ComparePass(h, p []byte) bool {
	hash, pass := []byte(h), []byte(p)
	err := bcrypt.CompareHashAndPassword(hash, pass)

	return err == nil
}