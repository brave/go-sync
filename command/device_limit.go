package command

import (
	"strings"
)

const (
	maxActiveDevices     int = 50
	highMaxActiveDevices int = 100
)

var (
	highDeviceLimitClientIDs = make(map[string]bool)
)

func LoadHighDeviceLimitClientIDs(clientIDList string) {
	highDeviceLimitClientIDs = make(map[string]bool)
	if clientIDList != "" {
		for id := range strings.SplitSeq(clientIDList, ",") {
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
