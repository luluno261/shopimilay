package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// GenerateID génère un ID unique
func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// NowTimestamp retourne le timestamp actuel en millisecondes
func NowTimestamp() int64 {
	return time.Now().UnixMilli()
}

// FormatDate formate une date
func FormatDate(t time.Time) string {
	return t.Format(time.RFC3339)
}

