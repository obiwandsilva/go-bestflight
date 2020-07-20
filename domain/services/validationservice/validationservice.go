package validationservice

import (
	"log"
	"regexp"
)

// IsValidAirport returns true if an airport name matches the expected format.
func IsValidAirport(airport string) bool {
	pattern := `^[A-Z]{3}$`
	match, err := regexp.MatchString(pattern, airport)
	if err != nil {
		log.Println(err)
	}

	return match
}
