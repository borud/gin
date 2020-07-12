package auth

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
)

func createRandomStateString() string {
	data := make([]byte, 128)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		log.Fatalf("Error generating random string: %v", err)
	}
	return base64.StdEncoding.EncodeToString(data)
}
