package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic("unable to generate secure ID: " + err.Error())
	}
	return hex.EncodeToString(bytes)[:length]
}
