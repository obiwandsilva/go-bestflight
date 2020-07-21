package routeservice

import (
	"container/heap"
	r "go-bestflight/domain/entities/routes"
	"go-bestflight/domain/errors"
	"strings"
)

type indexes map[interface{}]interface{}

type routesGraph [][]r.Connection

type mapper struct {
	indxs     indexes
	distances []int
	previous  []int
}

type dijkstraArgs struct {
	start int
	end   int
	dist  []int
	prev  []int
	indxs indexes
	g     routesGraph
}

const (
	maxInt = int(^uint(0) >> 1)
)

func convertRouteToNamed(route []int, indxs indexes) string {
	airports := []string{}

	for _, node := range route {
		airports = append(airports, indxs[node].(string))
	}

	return strings.Join(airports, " - ")
}

func buildMapper(airports []string) mapper {
	indxs := make(indexes)
	distances := make([]int, len(airports))
	previous := make([]int, len(airports))

	for i, airport := range airports {
		indxs[airport] = i
		indxs[i] = airport
		distances[i] = maxInt
		previous[i] = -1
	}

	return mapper{
		indxs:     indxs,
		distances: distances,
		previous:  previous,
	}
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
	prevNode := previous[end]
	route := append([]int{}, end, prevNode)

	for prevNode != start {
		prevNode = previous[prevNode]
		route = append(route, prevNode)
	}

	return reverseRoute(route)
}

// DijkstraSTP implements the Dijkstra's Shortest Path algorithm.
func DijkstraSTP(args dijkstraArgs) ([]int, int) {
	pq := NewPriorityQueue()
	visited := make([]bool, len(args.dist))

	args.dist[args.start] = 0
	args.prev[args.start] = args.start
	heap.Push(pq, &Item{
		node:     args.start,
		priority: 0,
	})

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
			args.prev[destinationNode] = nodeMinDistance.node

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

	if args.prev[args.end] == -1 {
		return []int{}, -1
	}

	bestRoute := reconstructRoute(args.start, args.end, args.prev)
	cost := args.dist[args.end]

	return bestRoute, cost
}

func findBestRoute(airports []string, routes r.Routes, boarding, destination string) (r.BestRoute, error) {
	m := buildMapper(airports)
	g := buildGraph(routes, m.indxs, len(m.distances))
	args := dijkstraArgs{
		start: m.indxs[boarding].(int),
		end:   m.indxs[destination].(int),
		dist:  m.distances,
		indxs: m.indxs,
		prev:  m.previous,
		g:     g,
	}
	bestRoute, cost := DijkstraSTP(args)

	if cost == maxInt || cost == -1 {
		return r.BestRoute{}, errors.NewBestRouteNotFoundErr()
	}

	best := r.BestRoute{
		Route: convertRouteToNamed(bestRoute, m.indxs),
		Cost:  cost,
	}

	return best, nil
}
