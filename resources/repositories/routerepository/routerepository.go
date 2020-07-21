package routerepository

import (
	r "go-bestflight/domain/entities/routes"
	"go-bestflight/resources/cache"
	"go-bestflight/resources/database"
	"go-bestflight/resources/file"
	"log"
)

// StoreRoute encapsulates the adding of new routes and airports to the database, cache and file.
func StoreRoute(route r.Route) error {
	database.StoreAirport(route.Boarding)
	database.StoreAirport(route.Destination)
	database.StoreRoute(route)

	err := file.Write(route)
	if err != nil {
		log.Printf("error when writing to file: %v", err)
		log.Println("removing route from database")
		database.DeleteRoute(route)
		return err
	}

	cache.AddRoute(route)

	return nil
}

// StoreRouteFromFile stores routes and airports from file into database and cache.
func StoreRouteFromFile(route r.Route) {
	database.StoreAirport(route.Boarding)
	database.StoreAirport(route.Destination)
	database.StoreRoute(route)
	cache.AddRoute(route)
}
