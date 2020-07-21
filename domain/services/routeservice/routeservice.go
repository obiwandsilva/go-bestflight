package routeservice

import (
	"errors"
	r "go-bestflight/domain/entities/routes"
	e "go-bestflight/domain/errors"
	validation "go-bestflight/domain/services/validationservice"
	"go-bestflight/resources/cache"
	"go-bestflight/resources/repositories/airportrepository"
	"go-bestflight/resources/repositories/routerepository"
	"log"
	"strings"
)

// AddNewRoute ...
func AddNewRoute(route r.Route) error {
	boarding := strings.ToUpper(route.Boarding)
	destination := strings.ToUpper(route.Destination)
	newRoute := r.Route{
		Boarding:    boarding,
		Destination: destination,
		Cost:        route.Cost,
	}

	if !validation.IsValidRoute(newRoute) {
		log.Printf("invalid route format: %v\n", newRoute)
		return e.NewInvalidRouteErr()
	}

	if routerepository.RouteExists(newRoute.Boarding, newRoute.Destination) {
		log.Printf("route already stored: %v\n", newRoute)
		return e.NewRouteAlreadyExistErr()
	}

	err := routerepository.StoreRoute(newRoute)
	if err != nil {
		return errors.New("could not create resource")
	}

	return nil
}

// LoadRoutes from file into database and cache.
func LoadRoutes(routes []r.Route) {
	for line, route := range routes {
		boarding := strings.ToUpper(route.Boarding)
		destination := strings.ToUpper(route.Destination)
		newRoute := r.Route{
			Boarding:    boarding,
			Destination: destination,
			Cost:        route.Cost,
		}

		if !validation.IsValidRoute(newRoute) {
			log.Printf("invalid format at line: %d\n", line)
			continue
		}

		if routerepository.RouteExists(newRoute.Boarding, newRoute.Destination) {
			log.Printf("route at line %d already stored: %v\n", line, newRoute)
			continue
		}

		routerepository.StoreRouteFromFile(newRoute)
	}
}

// GetBestRoute ...
func GetBestRoute(boarding string, destination string) (r.BestRoute, error) {
	board := strings.ToUpper(boarding)
	dest := strings.ToUpper(destination)

	if !validation.IsValidAirport(board) || !validation.IsValidAirport(dest) {
		return r.BestRoute{}, e.NewInvalidAirportErr("invalid")
	}

	if !airportrepository.IsRegistered(board) || !airportrepository.IsRegistered(dest) {
		return r.BestRoute{}, e.NewInvalidAirportErr("not registered")
	}

	airports := airportrepository.GetAllAirports()
	routes := cache.GetAllRoutes()

	bestRoute, err := findBestRoute(airports, routes, board, dest)
	if err != nil {
		log.Printf("error when getting best route for %s-%s: %v", board, dest, err)
		return r.BestRoute{}, err
	}

	return bestRoute, nil
}
