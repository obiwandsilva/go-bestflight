package cache

import (
	r "go-bestflight/domain/entities/routes"
	"sync"
)

// Memcache represents a cache service, e.g. Redis.
// It allows a constant ready to use easy routes format.
type Memcache struct {
	routes r.Routes
	sync.RWMutex
}

var (
	instance *Memcache
	once     sync.Once
)

// Connect iniciates the memcache instance only once.
func Connect() {
	once.Do(func() {
		instance = &Memcache{
			routes: make(r.Routes),
		}
	})
}

// Truncate ...
func Truncate() {
	instance = &Memcache{
		routes: make(r.Routes),
	}
}

// AddRoute ...
func AddRoute(route r.Route) r.Route {
	instance.Lock()
	defer instance.Unlock()

	destinations, ok := instance.routes[route.Boarding]
	if ok {
		dest := r.Connection{Airport: route.Destination, Cost: route.Cost}
		instance.routes[route.Boarding] = append(destinations, dest)
		return route
	}

	dest := []r.Connection{{Airport: route.Destination, Cost: route.Cost}}
	instance.routes[route.Boarding] = dest

	return route
}

// AddRoutes adds multiple routes to the cache.
func AddRoutes(routes []r.Route) {
	for _, route := range routes {
		AddRoute(route)
	}
}

// GetAllRoutes returna all current routes in cache.
func GetAllRoutes() r.Routes {
	routesCopy := make(r.Routes)

	instance.RLock()
	defer instance.RUnlock()

	for boarding, connections := range instance.routes {
		connectionsCopy := make([]r.Connection, len(connections))

		copy(connectionsCopy, connections)

		routesCopy[boarding] = connectionsCopy
	}

	return routesCopy
}
