package format

import (
	"bytes"
	"testing"
)

func TestGenerateRandomBytes(t *testing.T) {
	token, err := GenerateRandomBytes(0)
	if len(token) > 0 {
		t.Errorf("Token generated not having the adequate size, want '%d' got '%d'", 0, len(token))
	}
	if err != nil {
		t.Errorf("Got an error while generating token: %s", err.Error())
	}
	token, err = GenerateRandomBytes(32)
	if len(token) != 32 {
		t.Errorf("Token generated not having the adequate size, want '%d' got '%d'", 32, len(token))
	}
	if err != nil {
		t.Errorf("Got an error while generating token: %s", err.Error())
	}
	anotherToken, err := GenerateRandomBytes(32)
	if len(anotherToken) != 32 {
		t.Errorf("Token generated not having the adequate size, want '%d' got '%d'", 32, len(anotherToken))
	}
	if err != nil {
		t.Errorf("Got an error while generating token: %s", err.Error())
	}
	if bytes.Equal(token, anotherToken) {
		t.Errorf("The function doesn't return a randomized token, got '%s' twice", anotherToken)
	}
}

func TestGenerateRandomString(t *testing.T) {
	token, err := GenerateRandomString(0)
	if len(token) > 0 {
		t.Errorf("Token generated not having the adequate size, want '%d' got '%d'", 0, len(token))
	}
	if err != nil {
		t.Errorf("Got an error while generating token: %s", err.Error())
	}
	token, err = GenerateRandomString(32)
	if len(token) != 44 {
		t.Errorf("Token generated not having the adequate size, want '%d' got '%d'", 44, len(token))
	}
	if err != nil {
		t.Errorf("Got an error while generating token: %s", err.Error())
	}
	anotherToken, err := GenerateRandomString(32)
	if len(anotherToken) != 44 {
		t.Errorf("Token generated not having the adequate size, want '%d' got '%d'", 44, len(anotherToken))
	}
	if err != nil {
		t.Errorf("Got an error while generating token: %s", err.Error())
	}
	if token == anotherToken {
		t.Errorf("The function doesn't return a randomized token, got '%s' twice", anotherToken)
	}
}
