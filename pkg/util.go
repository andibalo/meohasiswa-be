package pkg

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"unicode/utf8"
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

func NullStrToStr(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func GenRandNumber(n int) string {

	minLimit := int(math.Pow10(n))
	maxLimit := int(math.Pow10(n - 1))
	randInt := int(rand.Float64() * float64(minLimit))
	if randInt < maxLimit {
		randInt += maxLimit
	}

	return strconv.Itoa(randInt)
}

func TruncateWithEllipsis(s string, maxLength int) string {
	if utf8.RuneCountInString(s) > maxLength {
		runes := []rune(s)
		return string(runes[:maxLength]) + "..."
	}
	return s
}

func ExtractDomainFromEmail(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email address: %s", email)
	}

	domain := parts[1]

	domainParts := strings.Split(domain, ".")
	if len(domainParts) < 2 {
		return "", fmt.Errorf("invalid domain in email: %s", email)
	}

	var mainDomain string
	if len(domainParts) > 2 && len(domainParts[len(domainParts)-2]) <= 3 {
		mainDomain = strings.Join(domainParts[len(domainParts)-3:], ".")
	} else {
		mainDomain = strings.Join(domainParts[len(domainParts)-2:], ".")
	}

	return mainDomain, nil
}
