package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GetSHA256(text string, salt string) string {
	hash := sha256.Sum256([]byte(text + salt))
	return hex.EncodeToString(hash[:])
}
