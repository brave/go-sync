package json

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
