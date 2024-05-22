package validating

import (
	"strings"

	"robpike.io/filter"
)

func ValidateCoordinates(coordinates string) bool {
	splitCoordinates := strings.Split(coordinates, ",")
	splitCoordinates[0] = strings.TrimSpace(splitCoordinates[0])
	splitCoordinates[1] = strings.TrimSpace(splitCoordinates[1])

	for _, coordinate := range splitCoordinates {
		if len(coordinate) < 2 || len(coordinate) > 8 {
			return false
		}

		sanitizedCoordinate := filter.Choose(coordinate, func(r rune) bool {
			return (r >= '0' && r <= '9') || r == '.'
		}).(string)

		if len(sanitizedCoordinate) != len(coordinate) {
			return false
		}
	}

	return true
}