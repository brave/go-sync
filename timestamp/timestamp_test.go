package timestamp_test

import (
	"encoding/json"
	"testing"

	"github.com/brave/go-sync/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestGetTimestamp(t *testing.T) {
	rsp, err := timestamp.GetTimestamp()
	assert.Nil(t, err)

	// Unmarshal to get the timestamp value
	time := timestamp.Timestamp{}
	err = json.Unmarshal(rsp, &time)
	assert.Nil(t, err)

	expectedJSON := "{\"timestamp\":\"" + time.Timestamp + "\"}"
	assert.Equal(t, expectedJSON, string(rsp))
}
