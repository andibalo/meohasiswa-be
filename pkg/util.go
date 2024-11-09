package pkg

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ToPointer[K any](val K) *K {
	return &val
}

func GetCursorData(cursor string) (string, string) {

	splitCursor := strings.Split(cursor, "_")

	return splitCursor[0], splitCursor[1]
}
