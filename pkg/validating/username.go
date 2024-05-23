// Purpose: This package is responsible for validating the input data from the user.
package validating

import "robpike.io/filter"

func ValidateUsername(username string) bool {
	if len(username) < 4 || len(username) > 16 {
		return false
	}

	sanitizedUsername := filter.Choose([]rune(username), func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	}).([]rune)

	return len(sanitizedUsername) == len(username)
}
