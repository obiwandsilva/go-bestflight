package database

import (
	"go-bestflight/domain/errors"
	"go-bestflight/domain/routes"
	"testing"

	"github.com/franela/goblin"
)

func TestDatabase(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Tests for Connect", func() {
		g.It("should always return the same instance of db", func() {
			Connect()

			instance.airportTable["GRU"] = struct{}{}
			instance.airportTable["SCL"] = struct{}{}

			Connect()

			_, ok := instance.airportTable["GRU"]
			_, ok2 := instance.airportTable["SCL"]

			g.Assert(ok).Equal(true)
			g.Assert(ok2).Equal(true)
		})
	})

	g.Describe("Tests for StoreRoute", func() {
		g.BeforeEach(func() {
			Connect()
		})

		g.AfterEach(func() {
			truncate()
		})

		g.It("should successfully store a route", func() {
			route := routes.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			result := StoreRoute(route)

			g.Assert(result).Equal(route)
		})

		g.It("should not return errors when storing the same route multiple times", func() {
			route1 := routes.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := routes.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}

			StoreRoute(route1)

			result := StoreRoute(route2)

			g.Assert(result).Equal(route2)
		})

		g.It("should successfully store a routes with equal boardings and different destinations", func() {
			route := routes.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := routes.Route{
				Boarding:    "GRU",
				Destination: "ORL",
				Cost:        56,
			}

			StoreRoute(route)

			result := StoreRoute(route2)

			g.Assert(result).Equal(route2)
		})

		g.It("should successfully store a routes with different boardings and same destinations", func() {
			route := routes.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := routes.Route{
				Boarding:    "ORL",
				Destination: "CDG",
				Cost:        5,
			}

			StoreRoute(route)

			result := StoreRoute(route2)

			g.Assert(result).Equal(route2)
		})

		g.It("should successfully store a routes with different boardings and destinations", func() {
			route := routes.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := routes.Route{
				Boarding:    "BRC",
				Destination: "SCL",
				Cost:        5,
			}

			StoreRoute(route)

			result := StoreRoute(route2)

			g.Assert(result).Equal(route2)
		})
	})

	g.Describe("Tests for GetRouteCost", func() {
		g.BeforeEach(func() {
			Connect()
		})

		g.AfterEach(func() {
			truncate()
		})

		g.It("should successfully return a cost for a stored route", func() {
			boarding := "GRU"
			destination := "ORL"
			cost := 56
			route := routes.Route{
				Boarding:    boarding,
				Destination: destination,
				Cost:        cost,
			}

			StoreRoute(route)

			result, err := GetRouteCost(boarding, destination)

			g.Assert(err).Equal(nil)
			g.Assert(result).Equal(cost)
		})

		g.It("should return an error for unexisting route", func() {
			boarding := "GRU"
			destination := "ORL"

			result, err := GetRouteCost(boarding, destination)

			g.Assert(err).Equal(errors.NewRouteNotFoundErr())
			g.Assert(result).Equal(-1)
		})

		g.It("should successfully return a cost for different stored routes", func() {
			route := routes.Route{
				Boarding:    "GRU",
				Destination: "BRC",
				Cost:        10,
			}
			route2 := routes.Route{
				Boarding:    "BRC",
				Destination: "SCL",
				Cost:        5,
			}
			route3 := routes.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route4 := routes.Route{
				Boarding:    "GRU",
				Destination: "SCL",
				Cost:        20,
			}
			route5 := routes.Route{
				Boarding:    "GRU",
				Destination: "ORL",
				Cost:        56,
			}
			route6 := routes.Route{
				Boarding:    "ORL",
				Destination: "CDG",
				Cost:        5,
			}
			route7 := routes.Route{
				Boarding:    "SCL",
				Destination: "ORL",
				Cost:        20,
			}

			StoreRoute(route)
			StoreRoute(route2)
			StoreRoute(route3)
			StoreRoute(route4)
			StoreRoute(route5)
			StoreRoute(route6)
			StoreRoute(route7)

			result, err := GetRouteCost("GRU", "BRC")
			result2, err2 := GetRouteCost("BRC", "SCL")
			result3, err3 := GetRouteCost("GRU", "CDG")
			result4, err4 := GetRouteCost("GRU", "SCL")
			result5, err5 := GetRouteCost("GRU", "ORL")
			result6, err6 := GetRouteCost("ORL", "CDG")
			result7, err7 := GetRouteCost("SCL", "ORL")
			result8, err8 := GetRouteCost("A", "B")

			g.Assert(err).Equal(nil)
			g.Assert(err2).Equal(nil)
			g.Assert(err3).Equal(nil)
			g.Assert(err4).Equal(nil)
			g.Assert(err5).Equal(nil)
			g.Assert(err6).Equal(nil)
			g.Assert(err7).Equal(nil)
			g.Assert(err8).Equal(errors.NewRouteNotFoundErr())
			g.Assert(result).Equal(10)
			g.Assert(result2).Equal(5)
			g.Assert(result3).Equal(75)
			g.Assert(result4).Equal(20)
			g.Assert(result5).Equal(56)
			g.Assert(result6).Equal(5)
			g.Assert(result7).Equal(20)
			g.Assert(result8).Equal(-1)
		})
	})

	g.Describe("Tests for StoreRoutes", func() {
		g.It("should successfully store multiple routes", func() {
			Connect()

			routes := []routes.Route{
				{
					Boarding:    "GRU",
					Destination: "BRC",
					Cost:        10,
				},
				{
					Boarding:    "BRC",
					Destination: "SCL",
					Cost:        5,
				},
				{
					Boarding:    "GRU",
					Destination: "CDG",
					Cost:        75,
				},
				{
					Boarding:    "GRU",
					Destination: "SCL",
					Cost:        20,
				},
				{
					Boarding:    "GRU",
					Destination: "ORL",
					Cost:        56,
				},
				{
					Boarding:    "ORL",
					Destination: "CDG",
					Cost:        5,
				},
				{
					Boarding:    "SCL",
					Destination: "ORL",
					Cost:        20,
				},
			}

			StoreRoutes(routes)

			result, err := GetRouteCost("GRU", "BRC")
			result2, err2 := GetRouteCost("BRC", "SCL")
			result3, err3 := GetRouteCost("GRU", "CDG")
			result4, err4 := GetRouteCost("GRU", "SCL")
			result5, err5 := GetRouteCost("GRU", "ORL")
			result6, err6 := GetRouteCost("ORL", "CDG")
			result7, err7 := GetRouteCost("SCL", "ORL")
			result8, err8 := GetRouteCost("A", "B")

			g.Assert(err).Equal(nil)
			g.Assert(err2).Equal(nil)
			g.Assert(err3).Equal(nil)
			g.Assert(err4).Equal(nil)
			g.Assert(err5).Equal(nil)
			g.Assert(err6).Equal(nil)
			g.Assert(err7).Equal(nil)
			g.Assert(err8).Equal(errors.NewRouteNotFoundErr())
			g.Assert(result).Equal(10)
			g.Assert(result2).Equal(5)
			g.Assert(result3).Equal(75)
			g.Assert(result4).Equal(20)
			g.Assert(result5).Equal(56)
			g.Assert(result6).Equal(5)
			g.Assert(result7).Equal(20)
			g.Assert(result8).Equal(-1)

			truncate()
		})
	})

	g.Describe("Tests for StoreAirport", func() {
		g.It("should successfully store an airport", func() {
			Connect()

			airport := "GRU"
			result := StoreAirport(airport)

			g.Assert(result).Equal(airport)

			truncate()
		})
	})

	g.Describe("Tests for GetAllAirports", func() {
		g.It("should successfully return all stored airports", func() {
			Connect()

			airports := map[string]struct{}{
				"GRU": {},
				"CDG": {},
				"ORL": {},
				"BRC": {},
			}

			for airport := range airports {
				StoreAirport(airport)
			}

			result := GetAllAirports()

			g.Assert(len(result)).Equal(4)

			for _, airport := range result {
				_, ok := airports[airport]

				g.Assert(ok).Equal(true)
			}

			truncate()
		})
	})
}
