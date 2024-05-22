package validating

import "robpike.io/filter"

func ValidateUsername(username string) bool {
	if len(username) < 4 || len(username) > 16 {
		return false
	}

	sanitizedUsername := filter.Choose(username, func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
	}).(string)

	return len(sanitizedUsername) == len(username)
}
