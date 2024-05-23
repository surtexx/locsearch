//go:build !integration
// +build !integration

package locationHistory

import "testing"

func TestGetDistanceTraveled(t *testing.T) {
	username := "test"
	startDate := "2024-01-01"
	endDate := "2024-12-30"

	distance := GetDistanceTraveled(username, startDate, endDate)
	if distance == 0 {
		t.Error("Expected distance to be greater than 0")
	}
}
