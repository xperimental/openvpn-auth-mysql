package main

import (
	"code.google.com/p/gopass"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"crypto/rand"
)

func main() {
	salt := createSalt(8)
	fmt.Printf("Salt: %s\n", salt)

	password, err := gopass.GetPass("Password: ")
	if err != nil {
		fmt.Printf("Error getting password: %v\n", err)
	}

	hash := sha256.Sum256([]byte(salt + password))
	hashString := hex.EncodeToString(hash[:])
	hashed := salt + "|sha256|" + hashString

	fmt.Printf("Password hash: %s\n", hashed)
}

func createSalt(length int) string {
	buffer := make([]byte, length / 2)
	_, err := rand.Read(buffer)
	if err != nil {
		fmt.Printf("Error creating salt: %v\n", err)
		return "error"
	}
	return hex.EncodeToString(buffer)
}
