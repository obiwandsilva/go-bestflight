package routeservice

import (
	"container/heap"
	"go-bestflight/domain/errors"
	r "go-bestflight/domain/routes"
	"strings"
)

type indexes map[interface{}]interface{}

type routesGraph [][]r.Connection

type dijkstraArgs struct {
	start int
	end   int
	dist  []int
	indxs indexes
	g     routesGraph
}

const (
	maxInt = int(^uint(0) >> 1)
)

func convertRouteToNamed(route []int, indxs indexes) string {
	var strBuilder strings.Builder

	for _, node := range route {
		strBuilder.WriteString(indxs[node].(string))
	}

	return strBuilder.String()
}

func findBestRoute(airports []string, routes r.Routes, boarding, destination string) (r.BestRoute, error) {
	indexes, distances := buildIndexesAndDistance(airports)
	g := buildGraph(routes, indexes, len(distances))
	args := dijkstraArgs{
		start: indexes[boarding].(int),
		end:   indexes[destination].(int),
		dist:  distances,
		indxs: indexes,
		g:     g,
	}
	bestRoute, cost := DijkstraSTP(args)

	if cost == maxInt {
		return r.BestRoute{}, errors.NewBestRouteNotFoundErr()
	}

	best := r.BestRoute{
		Route: convertRouteToNamed(bestRoute, indexes),
		Cost:  cost,
	}

	return best, nil
}

func buildIndexesAndDistance(airports []string) (indexes, []int) {
	indxs := make(indexes)
	distances := make([]int, len(airports))

	for i, airport := range airports {
		indxs[airport] = i
		indxs[i] = airport
		distances[i] = maxInt
	}

	return indxs, distances
}

func buildGraph(routes r.Routes, indxs indexes, graphSize int) routesGraph {
	graph := make([][]r.Connection, graphSize)

	for boarding, destinations := range routes {
		i := indxs[boarding].(int)
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

		for _, destination := range args.g[nodeMinDistance.node] {
			destinationNode := args.indxs[destination.Airport].(int)
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
