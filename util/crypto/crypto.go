package crypto

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"strings"
)

func GenerateMD5Hash(str string) (string, error) {
	str = strings.ToLower(strings.TrimSpace(str))
	hash := md5.New()
	_, err := hash.Write([]byte(str))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
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
