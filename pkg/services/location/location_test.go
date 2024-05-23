//go:build !integration
// +build !integration

package location

import (
	"testing"
)

func TestSearchUsers(t *testing.T) {
	coordinates := "39.13354,27.14438"
	radius := 10.0

	users := SearchUsers(coordinates, radius)
	if len(users) == 0 {
		t.Error("Expected users to be found")
	}
}
