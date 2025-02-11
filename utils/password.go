package utils

import(
	"crypto/sha256"
)


func Hash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return string(sum[:])
}