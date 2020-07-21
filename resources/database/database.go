package database

import (
	r "go-bestflight/domain/entities/routes"
	"go-bestflight/domain/errors"
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
func Connect() {
	once.Do(func() {
		instance = Database{
			routeTable:   make(map[string]map[string]int),
			airportTable: make(map[string]struct{}),
		}
	})
}

func Truncate() {
	instance = Database{
		routeTable:   make(map[string]map[string]int),
		airportTable: make(map[string]struct{}),
	}
}

// StoreRoute ...
func StoreRoute(route r.Route) r.Route {
	instance.Lock()
	defer instance.Unlock()

	_, okBoarding := instance.routeTable[route.Boarding]

	if okBoarding {
		instance.routeTable[route.Boarding][route.Destination] = route.Cost

		return route
	}

	instance.routeTable[route.Boarding] = map[string]int{
		route.Destination: route.Cost,
	}

	return route
}

// DeleteRoute deletes a given route from database.
func DeleteRoute(route r.Route) {
	instance.Lock()
	defer instance.Unlock()

	destinations, okBoarding := instance.routeTable[route.Boarding]

	if okBoarding {
		_, okDestination := destinations[route.Destination]

		if okDestination {
			delete(instance.routeTable[route.Boarding], route.Destination)
		}
	}

	if len(instance.routeTable[route.Boarding]) == 0 {
		delete(instance.routeTable, route.Boarding)
	}
}

// GetRouteCost ...
func GetRouteCost(boarding, destination string) (int, error) {
	instance.RLock()
	defer instance.RUnlock()

	connections, ok := instance.routeTable[boarding]
	if !ok {
		return -1, errors.NewRouteNotFoundErr()
	}

	cost, ok := connections[destination]
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

// GetAirport returns true if the specified airport is found.
func GetAirport(airport string) bool {
	_, ok := instance.airportTable[airport]

	return ok
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
