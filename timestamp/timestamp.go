package timestamp

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/brave/go-sync/utils"
)

// Timestamp is a structure used for timestamp responses.
type Timestamp struct {
	Timestamp string `json:"timestamp"`
}

// GetTimestamp returns the current timestamp in JSON format.
func GetTimestamp() (rsp []byte, err error) {
	time := Timestamp{Timestamp: strconv.FormatInt(utils.UnixMilli(time.Now()), 10)}
	rsp, err = json.Marshal(time)
	return
}
