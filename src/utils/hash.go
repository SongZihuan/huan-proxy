package utils

import (
	"crypto/sha256"
	"fmt"
)

func SHA256(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}
