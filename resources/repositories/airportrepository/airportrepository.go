package airportrepository

import "go-bestflight/resources/database"

// IsRegistered returns true if the specified airport exists.
func IsRegistered(airport string) bool {
	return database.GetAirport(airport)
}

// GetAllAirports returns all stored airports.
func GetAllAirports() []string {
	return database.GetAllAirports()
}
