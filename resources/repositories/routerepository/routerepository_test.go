package routerepository

import (
	r "go-bestflight/domain/entities/routes"
	"go-bestflight/resources/cache"
	"go-bestflight/resources/database"
	"go-bestflight/resources/file"
	"testing"

	"github.com/franela/goblin"
)

func TestRoutesRepository(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Tests for StoreRoute", func() {
		g.It("should successfully store a route on database, cache and file", func() {
			filePath := "test.csv"

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			route := r.Route{
				Boarding:    "XYZ",
				Destination: "ABC",
				Cost:        1000,
			}

			err := StoreRoute(route)

			g.Assert(err).Equal(nil)

			cost, _ := database.GetRouteCost(route.Boarding, route.Destination)
			bAirport := database.GetAirport(route.Boarding)
			dAirport := database.GetAirport(route.Destination)
			routesFromCache := cache.GetAllRoutes()
			routesFromFile, _ := file.ReadFile()

			g.Assert(cost).Equal(1000)
			g.Assert(bAirport).IsTrue()
			g.Assert(dAirport).IsTrue()
			g.Assert(len(routesFromCache[route.Boarding])).Equal(1)
			g.Assert(len(routesFromFile)).Equal(1)
			g.Assert(routesFromFile[0]).Equal(route)

			file.Remove()
		})

		g.It("should remove route from database if file writing fails", func() {
			filePath := "" // empty path will generate errors when opening file

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			route := r.Route{
				Boarding:    "XYZ",
				Destination: "ABC",
				Cost:        1000,
			}

			err := StoreRoute(route)

			g.Assert(err != nil).IsTrue()

			cost, _ := database.GetRouteCost(route.Boarding, route.Destination)
			bAirport := database.GetAirport(route.Boarding)
			dAirport := database.GetAirport(route.Destination)
			routesFromCache := cache.GetAllRoutes()

			g.Assert(cost).Equal(-1)
			g.Assert(bAirport).IsTrue()
			g.Assert(dAirport).IsTrue()
			g.Assert(len(routesFromCache[route.Boarding])).Equal(0)

			file.Remove()
		})
	})

	g.Describe("Tests for LoadRoutes", func() {
		g.It("should successfully load routes into the database and cache", func() {
			filePath := "test.csv"

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			routes := []r.Route{
				{
					Boarding:    "AAA",
					Destination: "BBB",
					Cost:        10,
				},
				{
					Boarding:    "AAA",
					Destination: "CCC",
					Cost:        20,
				},
				{
					Boarding:    "AAA",
					Destination: "DDD",
					Cost:        56,
				},
				{
					Boarding:    "AAA",
					Destination: "EEE",
					Cost:        75,
				},
				{
					Boarding:    "BBB",
					Destination: "CCC",
					Cost:        5,
				},
			}

			for _, r := range routes {
				file.Write(r)
			}

			LoadRoutes(routes)

			cost, _ := database.GetRouteCost(routes[0].Boarding, routes[0].Destination)
			cost2, _ := database.GetRouteCost(routes[1].Boarding, routes[1].Destination)
			cost3, _ := database.GetRouteCost(routes[2].Boarding, routes[2].Destination)
			cost4, _ := database.GetRouteCost(routes[3].Boarding, routes[3].Destination)
			cost5, _ := database.GetRouteCost(routes[4].Boarding, routes[4].Destination)
			airportA := database.GetAirport(routes[0].Boarding)
			airportB := database.GetAirport(routes[0].Destination)
			airportC := database.GetAirport(routes[1].Destination)
			airportD := database.GetAirport(routes[2].Destination)
			airportE := database.GetAirport(routes[3].Destination)
			routesFromCache := cache.GetAllRoutes()
			routesFromFile, _ := file.ReadFile()

			g.Assert(cost).Equal(10)
			g.Assert(cost2).Equal(20)
			g.Assert(cost3).Equal(56)
			g.Assert(cost4).Equal(75)
			g.Assert(cost5).Equal(5)

			g.Assert(airportA).IsTrue()
			g.Assert(airportB).IsTrue()
			g.Assert(airportC).IsTrue()
			g.Assert(airportD).IsTrue()
			g.Assert(airportE).IsTrue()

			g.Assert(len(routesFromCache[routes[0].Boarding])).Equal(4)
			g.Assert(len(routesFromCache[routes[4].Boarding])).Equal(1)

			g.Assert(len(routesFromFile)).Equal(5)

			g.Assert(routesFromFile[0]).Equal(routes[0])
			g.Assert(routesFromFile[1]).Equal(routes[1])
			g.Assert(routesFromFile[2]).Equal(routes[2])
			g.Assert(routesFromFile[3]).Equal(routes[3])
			g.Assert(routesFromFile[4]).Equal(routes[4])

			file.Remove()
		})
	})
}
