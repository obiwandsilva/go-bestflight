package database

import (
	"go-bestflight/domain/errors"
	"go-bestflight/domain/routes"
	"sync"
)

// Database is reponsible for storing routes and airports data in memory.
type Database struct {
	routeTable   map[string]map[string]int
	airportTable map[string]struct{}
	sync.RWMutex
}

var (
	db   Database
	once sync.Once
)

// Connect ...
func Connect() *Database {
	once.Do(func() {
		db = Database{
			routeTable:   make(map[string]map[string]int),
			airportTable: make(map[string]struct{}),
		}
	})

	return &db
}

func truncate() {
	db = Database{
		routeTable:   make(map[string]map[string]int),
		airportTable: make(map[string]struct{}),
	}
}

// StoreRoute ...
func StoreRoute(route routes.Route) routes.Route {
	db.Lock()
	defer db.Unlock()

	dest, okBoarding := db.routeTable[route.Boarding]

	if okBoarding {
		_, okDestination := dest[route.Destination]

		if okDestination {
			db.routeTable[route.Boarding][route.Destination] = route.Cost
			return route
		}

		db.routeTable[route.Boarding][route.Destination] = route.Cost

		return route
	}

	db.routeTable[route.Boarding] = map[string]int{
		route.Destination: route.Cost,
	}

	return route
}

// GetRouteCost ...
func GetRouteCost(boarding, destination string) (int, error) {
	db.RLock()
	defer db.RUnlock()

	destinations, ok := db.routeTable[boarding]
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
func StoreRoutes(routes []routes.Route) {
	for _, route := range routes {
		StoreRoute(route)
	}
}

// StoreAirport ...
func StoreAirport(airport string) string {
	db.Lock()
	defer db.Unlock()

	db.airportTable[airport] = struct{}{}

	return airport
}

// GetAllAirports ...
func GetAllAirports() []string {
	db.RLock()
	defer db.RUnlock()

	airports := []string{}

	for airport := range db.airportTable {
		airports = append(airports, airport)
	}

	return airports
}
