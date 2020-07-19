package routeservice

import (
	"container/heap"
	"go-bestflight/domain"
)

type index map[interface{}]interface{}

type dijkstraArgs struct {
	start   int
	end     int
	dist    []int
	indexes index
	graph   [][]domain.Destination
}

func findBestRoute(airports []string, routes domain.Routes, boarding, destination string) (string, int) {
	// indexes, distances := buildIndexAndDistance(airports)
	// graph := buildGraph(routes, indexes, len(distances))
	// boardingIndex := indexes[boarding]
	// destinationIndex := indexes[destination]

	return "", 0
}

func buildIndexesAndDistance(airports []string) (index, []int) {
	indexes := make(index)
	distances := make([]int, len(airports))
	maxInt := int(^uint(0) >> 1)

	for i, airport := range airports {
		indexes[airport] = i
		indexes[i] = airport
		distances[i] = maxInt
	}

	return indexes, distances
}

func buildGraph(routes domain.Routes, indexes index, graphSize int) [][]domain.Destination {
	graph := make([][]domain.Destination, graphSize)

	for boarding, destinations := range routes {
		i := indexes[boarding].(int)
		graph[i] = destinations
	}

	return graph
}

func reverseRoute(route []int) []int {
	length := len(route)

	if length < 2 {
		return route
	}

	for head, tail := 0, length-1; head < tail; head, tail = head+1, tail-1 {
		route[head], route[tail] = route[tail], route[head]
	}

	return route
}

func reconstructRoute(start, end int, previous []int) []int {
	route := []int{}
	prevNode := previous[end]
	route = append(route, end, prevNode)

	for prevNode != start {
		prevNode = previous[prevNode]
		route = append(route, prevNode)
	}

	return reverseRoute(route)
}

// DijkstraSTP implements the Dijkstra's Shortest Path algorithm.
func DijkstraSTP(args dijkstraArgs) ([]int, int) {
	pq := NewPriorityQueue()
	previous := make([]int, len(args.dist))
	visited := make([]bool, len(args.dist))

	heap.Push(pq, &Item{
		node:     args.start,
		priority: 0,
	})

	args.dist[args.start] = 0
	previous[args.start] = args.start

	for pq.Len() != 0 {
		var nodeMinDistance *Item = heap.Pop(pq).(*Item)

		if visited[nodeMinDistance.node] {
			continue
		}

		for _, destination := range args.graph[nodeMinDistance.node] {
			destinationNode := args.indexes[destination.Airport].(int)
			newDistance := args.dist[nodeMinDistance.node] + destination.Cost

			if newDistance >= args.dist[destinationNode] {
				continue
			}

			args.dist[destinationNode] = newDistance
			previous[destinationNode] = nodeMinDistance.node

			// Lazy implementation, but better than using an update on the current
			// PriorityQueue implementation.
			heap.Push(pq, &Item{
				node:     destinationNode,
				priority: destination.Cost,
			})
		}

		visited[nodeMinDistance.node] = true

		if nodeMinDistance.node == args.end {
			break
		}
	}

	bestRoute := reconstructRoute(args.start, args.end, previous)
	cost := args.dist[args.end]

	return bestRoute, cost
}
