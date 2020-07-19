package routeservice

import "go-bestflight/domain"

// DistanceInfo holds the information of the shortest distance from a starting point
// to the destination.
type DistanceInfo struct {
	Shortest int
	Previous string
}

type graph [][]int

func findBestRoute(routes domain.Routes, boarding, destination string) (string, int) {

}

func buildGraph(routes domain.Routes, indexes map[string]int) [][]int {
	graph := [][]int{}

	for boarding, dests := range routes {
		connections := []int{}
		
		for dest, cost := range dests {
			
		}
	}
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
func Dijkstra(board, dest string, airports []string) map[string]DistanceInfo {
	indexes, dist := buildAirportsIndexAndDistance(airports)
	visited := make([]bool, len(dist))
	pq := NewPriorityQueue()

	item := &Item{
		value:    board,
		priority: 0,
	}

	pq.Push(item)

	for pq.Len() != 0 {

	}
}
