package timestamp_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/brave/go-sync/auth"
	jsonschema "github.com/brave/go-sync/schema/json"
	"github.com/brave/go-sync/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestGetTimestamp(t *testing.T) {
	rsp, err := timestamp.GetTimestamp()
	assert.Nil(t, err)

	// Unmarshal to get the timestamp value
	timestampRsp := jsonschema.TimestampResponse{}
	err = json.Unmarshal(rsp, &timestampRsp)
	assert.Nil(t, err)

	expectedJSON := "{\"timestamp\":\"" + timestampRsp.Timestamp + "\",\"expires_in\":" + strconv.FormatInt(auth.TokenMaxDuration, 10) + "}"
	assert.Equal(t, expectedJSON, string(rsp))
}
