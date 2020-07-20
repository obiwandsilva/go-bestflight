package routeservice

import (
	"fmt"
	r "go-bestflight/domain/entities/routes"
	"testing"

	"github.com/franela/goblin"
)

func TestShortestPath(t *testing.T) {
	g := goblin.Goblin(t)

	airports := []string{
		"ORL",
		"BRC",
		"GRU",
		"CDG",
		"SCL",
	}

	routes := r.Routes{
		"GRU": []r.Connection{
			{Airport: "BRC", Cost: 10},
			{Airport: "CDG", Cost: 75},
			{Airport: "SCL", Cost: 20},
			{Airport: "ORL", Cost: 56},
		},
		"BRC": []r.Connection{
			{Airport: "SCL", Cost: 5},
		},
		"ORL": []r.Connection{
			{Airport: "CDG", Cost: 5},
		},
		"SCL": []r.Connection{
			{Airport: "ORL", Cost: 20},
		},
	}

	g.Describe("Tests for buildIndexAndDistance", func() {
		g.It("should successfully build a map of indxs and a slice of distances", func() {
			maxInt := int(^uint(0) >> 1)
			indxs, distances := buildIndexesAndDistance(airports)

			g.Assert(indxs[airports[0]]).Equal(0)
			g.Assert(indxs[airports[1]]).Equal(1)
			g.Assert(indxs[airports[2]]).Equal(2)
			g.Assert(indxs[airports[3]]).Equal(3)
			g.Assert(indxs[airports[4]]).Equal(4)

			g.Assert(len(distances)).Equal(len(airports))

			for _, distance := range distances {
				g.Assert(distance).Equal(maxInt)
			}
		})
	})

	g.Describe("Tests for buildGraph", func() {
		g.It("should successfully build a graph", func() {
			indxs, distances := buildIndexesAndDistance(airports)
			graph := buildGraph(routes, indxs, len(distances))

			g.Assert(len(graph)).Equal(len(distances))

			for _, airport := range airports {
				airportIndex := indxs[airport].(int)

				g.Assert(len(graph[airportIndex])).Equal(len(routes[airport]))
				g.Assert(graph[airportIndex]).Equal(routes[airport])
			}
		})
	})

	g.Describe("Tests for reverseRoute", func() {
		g.It("should successfully reverse a slice", func() {
			routeOne := []int{5}
			routeTwo := []int{5, 6}
			routeThree := []int{5, 6, 7}
			routeMany := []int{5, 3, 7, 9, 10, 35, 0}

			g.Assert(reverseRoute(routeOne)).Equal([]int{5})
			g.Assert(reverseRoute(routeTwo)).Equal([]int{6, 5})
			g.Assert(reverseRoute(routeThree)).Equal([]int{7, 6, 5})
			g.Assert(reverseRoute(routeMany)).Equal(
				[]int{0, 35, 10, 9, 7, 3, 5},
			)
		})
	})

	g.Describe("Tests for reconstructRoute", func() {
		g.It("should successfully reconstruct a route", func() {
			// Considering routes
			// 	GRU,BRC,10
			// 	BRC,SCL,5
			// 	GRU,CDG,75
			// 	GRU,SCL,20
			// 	GRU,ORL,56
			// 	ORL,CDG,5
			// 	SCL,ORL,20
			// Where indxs are:
			// 	GRU = 0
			// 	BRC = 1
			// 	SCL = 2
			// 	ORL = 3
			// 	CDG = 4

			start := 0
			end := 4
			previous := []int{0, 0, 1, 2, 3}

			result := reconstructRoute(start, end, previous)

			g.Assert(result).Equal(reverseRoute([]int{4, 3, 2, 1, 0}))
		})

		g.It("should successfully reconstruct a route to close node", func() {
			// Considering routes
			// 	GRU,BRC,10
			// 	BRC,SCL,5
			// 	GRU,CDG,75
			// 	GRU,SCL,20
			// 	GRU,ORL,56
			// 	ORL,CDG,5
			// 	SCL,ORL,20
			// Where indxs are:
			// 	GRU = 0
			// 	BRC = 1
			// 	SCL = 2
			// 	ORL = 3
			// 	CDG = 4

			start := 0
			end := 1
			previous := []int{0, 0, 1, 2, 3}

			result := reconstructRoute(start, end, previous)

			g.Assert(result).Equal(reverseRoute([]int{1, 0}))
		})
	})

	g.Describe("Tests for DijkstraSTP", func() {
		g.It("should retrieve the shortest path for a short distance", func() {
			indxs, distances := buildIndexesAndDistance(airports)
			gph := buildGraph(routes, indxs, len(distances))
			args := dijkstraArgs{
				start: indxs["BRC"].(int),
				end:   indxs["SCL"].(int),
				dist:  distances,
				indxs: indxs,
				g:     gph,
			}
			bestRoute, cost := DijkstraSTP(args)

			expectedRoute := []int{
				indxs["BRC"].(int),
				indxs["SCL"].(int),
			}

			g.Assert(bestRoute).Equal(expectedRoute)
			g.Assert(cost).Equal(5)
		})

		g.It("should retrieve the shortest path for a long distance", func() {
			indxs, distances := buildIndexesAndDistance(airports)
			graph := buildGraph(routes, indxs, len(distances))
			args := dijkstraArgs{
				start: indxs["GRU"].(int),
				end:   indxs["CDG"].(int),
				dist:  distances,
				indxs: indxs,
				g:     graph,
			}
			bestRoute, cost := DijkstraSTP(args)

			expectedRoute := []int{
				indxs["GRU"].(int),
				indxs["BRC"].(int),
				indxs["SCL"].(int),
				indxs["ORL"].(int),
				indxs["CDG"].(int),
			}

			g.Assert(bestRoute).Equal(expectedRoute)
			g.Assert(cost).Equal(40)
		})

		g.It("should retrieve the shortest path after add new routes", func() {
			newRoutes := make(r.Routes)
			for k, v := range routes {
				newRoutes[k] = v
			}
			newRoutes["X"] = []r.Connection{
				{Airport: "GRU", Cost: 7},
			}
			newAirports := append(airports, "X")

			indxs, distances := buildIndexesAndDistance(newAirports)
			graph := buildGraph(newRoutes, indxs, len(distances))

			args := dijkstraArgs{
				start: indxs["X"].(int),
				end:   indxs["GRU"].(int),
				dist:  distances,
				indxs: indxs,
				g:     graph,
			}

			bestRoute, cost := DijkstraSTP(args)

			expectedRoute := []int{
				indxs["X"].(int),
				indxs["GRU"].(int),
			}

			g.Assert(bestRoute).Equal(expectedRoute)
			g.Assert(cost).Equal(7)
		})

		g.It("should retrieve the shortest path even with reversed routes", func() {
			newRoutes := make(r.Routes)
			for k, v := range routes {
				newRoutes[k] = v
			}
			newRoutes["CDG"] = []r.Connection{
				{Airport: "GRU", Cost: 70},
			}
			newRoutes["ORL"] = append(newRoutes["ORL"], r.Connection{Airport: "GRU", Cost: 50})

			indxs, distances := buildIndexesAndDistance(airports)
			graph := buildGraph(newRoutes, indxs, len(distances))

			args := dijkstraArgs{
				start: indxs["GRU"].(int),
				end:   indxs["CDG"].(int),
				dist:  distances,
				indxs: indxs,
				g:     graph,
			}

			bestRoute, cost := DijkstraSTP(args)

			expectedRoute := []int{
				indxs["GRU"].(int),
				indxs["BRC"].(int),
				indxs["SCL"].(int),
				indxs["ORL"].(int),
				indxs["CDG"].(int),
			}

			g.Assert(bestRoute).Equal(expectedRoute)
			g.Assert(cost).Equal(40)
		})

		g.It("should retrieve the shortest path after add new routes disassociated from others", func() {
			newRoutes := make(r.Routes)
			for k, v := range routes {
				newRoutes[k] = v
			}
			newRoutes["X"] = []r.Connection{
				{Airport: "Y", Cost: 15},
			}
			newRoutes["Y"] = []r.Connection{
				{Airport: "X", Cost: 15},
				{Airport: "Z", Cost: 16},
			}
			newAirports := append(airports, "X", "Y", "Z")

			indxs, distances := buildIndexesAndDistance(newAirports)
			graph := buildGraph(newRoutes, indxs, len(distances))

			args := dijkstraArgs{
				start: indxs["X"].(int),
				end:   indxs["Z"].(int),
				dist:  distances,
				indxs: indxs,
				g:     graph,
			}

			bestRoute, cost := DijkstraSTP(args)

			expectedRoute := []int{
				indxs["X"].(int),
				indxs["Y"].(int),
				indxs["Z"].(int),
			}

			g.Assert(bestRoute).Equal(expectedRoute)
			g.Assert(cost).Equal(31)
		})

		g.It("should retrieve the shortest path after add new routes disassociated from others", func() {
			newRoutes := make(r.Routes)
			for k, v := range routes {
				newRoutes[k] = v
			}
			newRoutes["X"] = []r.Connection{
				{Airport: "Y", Cost: 15},
			}
			newAirports := append(airports, "X", "Y")

			indxs, distances := buildIndexesAndDistance(newAirports)
			graph := buildGraph(newRoutes, indxs, len(distances))

			args := dijkstraArgs{
				start: indxs["ORL"].(int),
				end:   indxs["X"].(int),
				dist:  distances,
				indxs: indxs,
				g:     graph,
			}

			bestRoute, cost := DijkstraSTP(args)

			fmt.Println(indxs, bestRoute, cost)

			// expectedRoute := []int{
			// 	indxs["X"].(int),
			// 	indxs["Y"].(int),
			// 	indxs["Z"].(int),
			// }

			// g.Assert(bestRoute).Equal(expectedRoute)
			// g.Assert(cost).Equal(31)
		})
	})
}
