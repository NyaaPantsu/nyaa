package crypto

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func GenerateMD5Hash(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	hash := md5.New()
	hash.Write([]byte(email))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
func GenerateRandomToken16() (string, error) {
	return GenerateRandomToken(16)
}

func GenerateRandomToken32() (string, error) {
	return GenerateRandomToken(32)
}

func GenerateRandomToken(n int) (string, error) {
	token := make([]byte, n)
	_, err := rand.Read(token)
	// %x	base 16, lower-case, two characters per byte
	return fmt.Sprintf("%x", token), err

}
