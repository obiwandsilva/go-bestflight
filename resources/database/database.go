package database

import (
	"go-bestflight/domain/errors"
	r "go-bestflight/domain/routes"
	"sync"
)

// Database is reponsible for storing routes and airports data in memory.
type Database struct {
	routeTable   map[string]map[string]int
	airportTable map[string]struct{}
	sync.RWMutex
}

var (
	instance Database
	once     sync.Once
)

// Connect ...
func Connect() *Database {
	once.Do(func() {
		instance = Database{
			routeTable:   make(map[string]map[string]int),
			airportTable: make(map[string]struct{}),
		}
	})

	return &instance
}

func truncate() {
	instance = Database{
		routeTable:   make(map[string]map[string]int),
		airportTable: make(map[string]struct{}),
	}
}

// StoreRoute ...
func StoreRoute(route r.Route) r.Route {
	instance.Lock()
	defer instance.Unlock()

	dest, okBoarding := instance.routeTable[route.Boarding]

	if okBoarding {
		_, okDestination := dest[route.Destination]

		if okDestination {
			instance.routeTable[route.Boarding][route.Destination] = route.Cost
			return route
		}

		instance.routeTable[route.Boarding][route.Destination] = route.Cost

		return route
	}

	instance.routeTable[route.Boarding] = map[string]int{
		route.Destination: route.Cost,
	}

	return route
}

// GetRouteCost ...
func GetRouteCost(boarding, destination string) (int, error) {
	instance.RLock()
	defer instance.RUnlock()

	destinations, ok := instance.routeTable[boarding]
	if !ok {
		return -1, errors.NewRouteNotFoundErr()
	}

	cost, ok := destinations[destination]
	if !ok {
		return -1, errors.NewRouteNotFoundErr()
	}

	return cost, nil
}

// StoreRoutes ...
func StoreRoutes(routes []r.Route) {
	for _, route := range routes {
		StoreRoute(route)
	}
}

// StoreAirport ...
func StoreAirport(airport string) string {
	instance.Lock()
	defer instance.Unlock()

	instance.airportTable[airport] = struct{}{}

	return airport
}

// GetAllAirports ...
func GetAllAirports() []string {
	instance.RLock()
	defer instance.RUnlock()

	airports := []string{}

	for airport := range instance.airportTable {
		airports = append(airports, airport)
	}

	return airports
}
