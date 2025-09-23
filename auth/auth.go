package auth

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// TokenMaxDuration specifies the lifetime for each access token.
	TokenMaxDuration int64  = 86400 * 1e3 // Milliseconds
	bearerPrefix     string = "Bearer "
	tokenRE          string = `^(?P<TimestampHex>[a-fA-F0-9]+)\|(?P<SignedTimestampHex>[a-fA-F0-9]+)\|(?P<PublicKeyHex>[a-fA-F0-9]+)$`
)

// Token represents the values we have in access tokens.
type Token struct {
	TimestampHex       string
	SignedTimestampHex string
	PublicKeyHex       string
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// authenticate parses the access token and verifies the signature and the
// timestamp is within +/- 1 day. The acccess token format is:
// base64_encode(timestamp_hex_encoded|signed_timestamp_hex_encoded|public_key_hex_encoded)
func authenticate(tkn string) (string, error) {
	base64DecodedBytes, err := base64.URLEncoding.DecodeString(tkn)
	if err != nil {
		return "", fmt.Errorf("error doing base64 decoding: %w", err)
	}

	re := regexp.MustCompile(tokenRE)
	m := re.FindStringSubmatch(string(base64DecodedBytes))
	if m == nil {
		return "", fmt.Errorf("invalid token format")
	}
	token := Token{TimestampHex: m[1], SignedTimestampHex: m[2], PublicKeyHex: m[3]}

	// Verify the signature.
	publicKey, err := hex.DecodeString(token.PublicKeyHex)
	if err != nil {
		return "", fmt.Errorf("error decoding hex string: %w", err)
	}
	timestampBytes, err := hex.DecodeString(token.TimestampHex)
	if err != nil {
		return "", fmt.Errorf("error decoding hex string: %w", err)
	}
	signedTimestamp, err := hex.DecodeString(token.SignedTimestampHex)
	if err != nil {
		return "", fmt.Errorf("error decoding hex string: %w", err)
	}
	if !ed25519.Verify(publicKey, timestampBytes, signedTimestamp) {
		return "", fmt.Errorf("signature verification failed")
	}

	var timestamp int64
	timestamp, err = strconv.ParseInt(string(timestampBytes), 10, 64)
	if err != nil {
		return "", fmt.Errorf("error parsing timestamp: %w", err)
	}

	// Verify that this token is within +/- 1 day.
	if abs(time.Now().UnixMilli()-timestamp) > TokenMaxDuration {
		return "", fmt.Errorf("token is expired")
	}

	blockedIDs := strings.Split(os.Getenv("BLOCKED_CLIENT_IDS"), ",")
	for _, id := range blockedIDs {
		if token.PublicKeyHex == id {
			return "", fmt.Errorf("This client ID is blocked")
		}
	}

	return token.PublicKeyHex, nil
}

// Authorize extracts the authorize token from the HTTP request and verify if
// the token is valid or not. It returns the clientID if the token is valid,
// otherwise, an empty string will be returned.
func Authorize(r *http.Request) (string, error) {
	var token string
	// Extract token from the header.
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		if !strings.HasPrefix(token, bearerPrefix) {
			return "", fmt.Errorf("Not a valid token")
		}
		token = strings.TrimPrefix(token, bearerPrefix)
	}
	if token == "" {
		return "", fmt.Errorf("Not a valid token")
	}

	// Verify token
	clientID, err := authenticate(token)
	if err != nil {
		return "", fmt.Errorf("error authorizing: %w", err)
	}
	return clientID, nil
}
