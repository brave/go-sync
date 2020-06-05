package json

// TimestampResponse is a structure used for timestamp responses.
type TimestampResponse struct {
	// Client will compose the actual token from it, it's named access_token to
	// avoid patching into chromium on client side.
	Timestamp string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}
