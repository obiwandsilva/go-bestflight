package validationservice

import (
	r "go-bestflight/domain/entities/routes"
	"log"
	"regexp"
)

const (
	min = 1
	max = 1000000
)

func isValidCost(cost int) bool {
	return (cost >= min) && (cost <= max)
}

func IsValidAirport(airport string) bool {
	pattern := `^[A-Z]{3}$`
	match, err := regexp.MatchString(pattern, airport)
	if err != nil {
		log.Println(err)
	}

	return match
}

// IsValidRoute ...
func IsValidRoute(route r.Route) bool {
	return IsValidAirport(route.Boarding) &&
		IsValidAirport(route.Destination) &&
		isValidCost(route.Cost)
}
