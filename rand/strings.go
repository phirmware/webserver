package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// RememberTokenBytes is the integer for the number of bytes wanted
const RememberTokenBytes = 32

// Bytes help generate n random bytes
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will take in n and use Bytes to generate n number of bytes the covert to a string and return
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken generates tokens of predetermined byte slice
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}
