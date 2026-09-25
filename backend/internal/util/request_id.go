package util

import (
	"crypto/rand"
	"encoding/hex"
)

// NewRequestID 随机请求 ID。
func NewRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
