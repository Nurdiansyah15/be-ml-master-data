package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateUniqueFileName menghasilkan nama file yang unik berdasarkan tipe dan waktu
func GenerateUniqueFileName(entityType string) string {
	// Menghasilkan random string
	randomBytes := make([]byte, 8) // 8 byte untuk random string
	_, _ = rand.Read(randomBytes)
	randomString := hex.EncodeToString(randomBytes)

	// Mendapatkan timestamp
	timestamp := time.Now().Unix()

	// Menghasilkan nama file dengan format yang diinginkan
	return fmt.Sprintf("%s_%s_%d", entityType, randomString, timestamp)
}
