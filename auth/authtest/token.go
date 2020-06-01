package authtest

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
)

// GenerateToken generates token from a given timestamp for tests to use.
func GenerateToken(timestamp int64) (string, string, string, error) {
	timestampBytes := []byte(strconv.FormatInt(timestamp, 10))
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", "", fmt.Errorf("generate key error: %w", err)
	}

	signedTimestampBytes := ed25519.Sign(privateKey, timestampBytes)
	publicKeyHex := hex.EncodeToString(publicKey)

	token := hex.EncodeToString(timestampBytes) + "|" +
		hex.EncodeToString(signedTimestampBytes) + "|" + publicKeyHex

	return base64.URLEncoding.EncodeToString([]byte(token)), token, publicKeyHex, nil
}
