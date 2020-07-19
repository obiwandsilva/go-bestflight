package routeservice

import "go-bestflight/domain"

// DistanceInfo holds the information of the shortest distance from a starting point
// to the destination.
type DistanceInfo struct {
	Shortest int
	Previous string
}

func findBestRoute(routes domain.Routes, boarding, destination string) (string, int) {
	return "", 0
}

func buildGraph(routes domain.Routes, indexes map[string]int) {
}

func buildAirportsIndexAndDistance(airports []string) (map[string]int, []int) {
	indexes := map[string]int{}
	distances := []int{}

	for i, airport := range airports {
		indexes[airport] = i
		distances[i] = int(^uint(0) >> 1)
	}

	return indexes, distances
}

// Dijkstra implements the Dijkstra's Shortest Path algorithm.
func Dijkstra(board, dest int, airports []string) map[string]DistanceInfo {
	return nil
}
