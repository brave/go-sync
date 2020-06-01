package json

// TimestampResponse is a structure used for timestamp responses.
type TimestampResponse struct {
	Timestamp string `json:"timestamp"`
	ExpiresIn int64  `json:"expires_in"`
}
