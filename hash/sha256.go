package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

// SHA256 使用 SHA256 计算字符串 Hash
func SHA256(v string) string {
	h := sha256.New()
	h.Write([]byte(v))
	return hex.EncodeToString(h.Sum(nil))
}
