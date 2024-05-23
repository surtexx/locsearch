// Purpose: This package is responsible for validating the input data from the user.
package validating

import (
	"strings"
)

// ISO 8601 format (YYY-MM-DDYTHH:MM:SS+HH:MM)
func ValiDate(date string) bool {
	dateSplit := strings.Split(date, "T")
	date = dateSplit[0]
	time := dateSplit[1]

	timeSplit := strings.Split(time, " ")
	time = timeSplit[0]
	timeZone := timeSplit[1]

	dateSplit = strings.Split(date, "-")
	year := dateSplit[0]
	month := dateSplit[1]
	day := dateSplit[2]

	if len(year) != 4 || len(month) != 2 || len(day) != 2 {
		return false
	}

	if len(timeZone) != 5 {
		return false
	}

	if len(time) != 8 {
		return false
	}

	return true
}
