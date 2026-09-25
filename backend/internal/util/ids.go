package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// NewID 生成短实体 ID（如 art-3f9a2c）。
func NewID(prefix string) string {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	}
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(b))
}
