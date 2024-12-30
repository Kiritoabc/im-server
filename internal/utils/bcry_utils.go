package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashPassword 使用 SHA-256 对密码进行加密
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// ComparePassword 对比输入的密码和哈希值
func ComparePassword(hashedPassword, password string) bool {
	return HashPassword(password) == hashedPassword
}
