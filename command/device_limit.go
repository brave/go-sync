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
	highDeviceLimitClientIDs map[string]bool
)

func init() {
	clientIDsEnv := os.Getenv("HIGH_DEVICE_LIMIT_CLIENT_IDS")
	LoadHighDeviceLimitClientIDs(clientIDsEnv)
}

func LoadHighDeviceLimitClientIDs(clientIDList string) {
	highDeviceLimitClientIDs = make(map[string]bool)
	if clientIDList != "" {
		ids := strings.Split(clientIDList, ",")
		for _, id := range ids {
			highDeviceLimitClientIDs[strings.ToLower(strings.TrimSpace(id))] = true
		}
	}
}

func hasReachedDeviceLimit(activeDevices int, clientID string) bool {
	limit := maxActiveDevices
	if highDeviceLimitClientIDs[strings.ToLower(clientID)] {
		limit = highMaxActiveDevices
	}
	return activeDevices >= limit
}
