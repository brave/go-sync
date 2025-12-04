package command

import (
	"os"
	"strings"
)

const (
	maxActiveDevices     int = 50
	highMaxActiveDevices int = 100
)

var (
	HighDeviceLimitClientIDs map[string]bool
)

func init() {
	clientIDsEnv := os.Getenv("HIGH_DEVICE_LIMIT_CLIENT_IDS")
	LoadHighDeviceLimitClientIDs(clientIDsEnv)
}

func LoadHighDeviceLimitClientIDs(clientIDList string) {
	HighDeviceLimitClientIDs = make(map[string]bool)
	if clientIDList != "" {
		ids := strings.Split(clientIDList, ",")
		for _, id := range ids {
			HighDeviceLimitClientIDs[strings.ToLower(strings.TrimSpace(id))] = true
		}
	}
}

func checkDeviceLimit(activeDevices int, clientID string) bool {
	limit := maxActiveDevices
	if HighDeviceLimitClientIDs[strings.ToLower(clientID)] {
		limit = highMaxActiveDevices
	}
	return activeDevices >= limit
}
