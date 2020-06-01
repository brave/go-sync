package timestamp

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/brave/go-sync/auth"
	jsonschema "github.com/brave/go-sync/schema/json"
	"github.com/brave/go-sync/utils"
)

// GetTimestamp returns the current timestamp in JSON format.
func GetTimestamp() (rsp []byte, err error) {
	time := jsonschema.TimestampResponse{
		Timestamp: strconv.FormatInt(utils.UnixMilli(time.Now()), 10),
		ExpiresIn: auth.TokenMaxDuration,
	}
	rsp, err = json.Marshal(time)
	return
}
