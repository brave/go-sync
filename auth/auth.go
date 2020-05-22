package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brave/go-sync/datastore"
	"github.com/brave/go-sync/utils"
)

var (
	timestampMaxDuration int64 = 120 * 1e3 // Milliseconds, modifiable in tests.
)

const (
	tokenMaxDuration int64  = 86400 * 1e3 // Milliseconds
	bearerPrefix     string = "Bearer "
	nRandBytes       int    = 32 // Number of random bytes used to generate the access token.
)

// Request is a struct used for authenication requests.
type Request struct {
	PublicKey       string `json:"public_key"`
	Timestamp       string `json:"timestamp"`
	SignedTimestamp string `json:"signed_timestamp"`
}

// Response is a struct used for authenication responses.
type Response struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// generateToken generates n random bytes and encoded it as base64 string.
func generateToken(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Authenticate validates client's auth requests and provides the reply.
func Authenticate(r *http.Request, db datastore.Datastore) (string, []byte, error) {
	var rsp []byte

	err := r.ParseForm()
	if err != nil {
		return "", nil, err
	}
	req := &Request{
		PublicKey:       r.PostFormValue("client_id"),
		Timestamp:       r.PostFormValue("timestamp"),
		SignedTimestamp: r.PostFormValue("client_secret"),
	}

	// Verify the signature.
	publicKey, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		return "", nil, err
	}
	timestampBytes, err := hex.DecodeString(req.Timestamp)
	if err != nil {
		return "", nil, err
	}
	signedTimestamp, err := hex.DecodeString(req.SignedTimestamp)
	if err != nil {
		return "", nil, err
	}
	if !ed25519.Verify(publicKey, timestampBytes, signedTimestamp) {
		return "", nil, fmt.Errorf("signature verification failed")
	}

	var timestamp int64
	timestamp, err = strconv.ParseInt(string(timestampBytes), 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("parse timestamp error")
	}

	// Verify the timestamp is not outdated
	if utils.UnixMilli(time.Now())-timestamp > timestampMaxDuration {
		return "", nil, fmt.Errorf("timestamp is outdated")
	}

	// Create a new token, save it in DB, and return it.
	expireAt := utils.UnixMilli(time.Now().Add(time.Duration(tokenMaxDuration) * time.Millisecond))
	token, err := generateToken(nRandBytes)
	if err != nil {
		return "", nil, err
	}
	err = db.InsertClientToken(req.PublicKey, token, expireAt)
	if err != nil {
		return "", nil, err
	}

	authRsp := Response{AccessToken: token, ExpiresIn: tokenMaxDuration}
	rsp, err = json.Marshal(authRsp)
	return token, rsp, err
}

// Authorize extracts the authorize token from the HTTP request and query the
// database to return the clientID associated with that token if the token is
// valid, otherwise, an empty string will be returned.
func Authorize(db datastore.Datastore, r *http.Request) (string, error) {
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

	// Query clients table for the token to return the clientID.
	return db.GetClientID(token)
}
