package routeservice

import (
	r "go-bestflight/domain/entities/routes"
	"go-bestflight/domain/errors"
	"go-bestflight/resources/cache"
	"go-bestflight/resources/database"
	"go-bestflight/resources/file"
	"strings"
	"testing"

	"github.com/franela/goblin"
)

func TestRouteService(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Tests for AddNewRoute", func() {
		g.It("should insert a new route applying business logic", func() {
			filePath := "test.csv"
			defer file.Remove()

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

			err := AddNewRoute(route)

			g.Assert(err).Equal(nil)

			cost, _ := database.GetRouteCost(route.Boarding, route.Destination)
			airport1 := database.GetAirport(route.Boarding)
			airport2 := database.GetAirport(route.Destination)
			routesFromCache := cache.GetAllRoutes()
			routesFromFile, _ := file.ReadFile()

			g.Assert(cost).Equal(1000)
			g.Assert(airport1).IsTrue()
			g.Assert(airport2).IsTrue()
			g.Assert(len(routesFromCache[route.Boarding])).Equal(1)
			g.Assert(len(routesFromFile)).Equal(1)
			g.Assert(routesFromFile[0]).Equal(route)
		})

		g.It("should insert a new route even with lowercase airports", func() {
			filePath := "test.csv"
			defer file.Remove()

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			route := r.Route{
				Boarding:    "xyz",
				Destination: "abc",
				Cost:        1000,
			}
			boarding := strings.ToUpper(route.Boarding)
			destination := strings.ToUpper(route.Destination)

			err := AddNewRoute(route)

			g.Assert(err).Equal(nil)

			cost, _ := database.GetRouteCost(boarding, destination)
			airport1 := database.GetAirport(boarding)
			airport2 := database.GetAirport(destination)
			routesFromCache := cache.GetAllRoutes()
			routesFromFile, _ := file.ReadFile()

			g.Assert(cost).Equal(1000)
			g.Assert(airport1).IsTrue()
			g.Assert(airport2).IsTrue()
			g.Assert(len(routesFromCache[boarding])).Equal(1)
			g.Assert(len(routesFromFile)).Equal(1)
			g.Assert(routesFromFile[0].Boarding).Equal(boarding)
			g.Assert(routesFromFile[0].Destination).Equal(destination)
			g.Assert(routesFromFile[0].Cost).Equal(route.Cost)
		})

		g.It("should return a InvalidRouteErr for invalid route format", func() {
			route := r.Route{
				Boarding:    "XY",
				Destination: "ABC",
				Cost:        1000,
			}

			err := AddNewRoute(route)

			g.Assert(err).Equal(errors.NewInvalidRouteErr())
		})

		g.It("should return a NewRouteAlreadyExistErr for an already stored route", func() {
			filePath := "test.csv"
			defer file.Remove()

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

			err := AddNewRoute(route)

			g.Assert(err).Equal(nil)

			err = AddNewRoute(route)

			g.Assert(err).Equal(errors.NewRouteAlreadyExistErr())
		})
	})

	g.Describe("Tests for LoadRoutes", func() {
		g.It("should load new routes into database and cache", func() {
			filePath := "test.csv"
			defer file.Remove()

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			routes := []r.Route{
				{
					Boarding:    "XYZ",
					Destination: "ABC",
					Cost:        1000,
				},
				{
					Boarding:    "XYZ",
					Destination: "EFG",
					Cost:        10,
				},
			}

			LoadRoutes(routes)

			cost, _ := database.GetRouteCost(routes[0].Boarding, routes[0].Destination)
			cost2, _ := database.GetRouteCost(routes[1].Boarding, routes[1].Destination)
			airport1 := database.GetAirport(routes[0].Boarding)
			airport2 := database.GetAirport(routes[0].Destination)
			airport3 := database.GetAirport(routes[1].Destination)
			routesFromCache := cache.GetAllRoutes()
			routesFromFile, _ := file.ReadFile()

			g.Assert(cost).Equal(1000)
			g.Assert(cost2).Equal(10)
			g.Assert(airport1).IsTrue()
			g.Assert(airport2).IsTrue()
			g.Assert(airport3).IsTrue()
			g.Assert(len(routesFromCache[routes[0].Boarding])).Equal(2)
			g.Assert(len(routesFromFile)).Equal(0)
		})

		g.It("should skip invalid and duplicate routes", func() {
			filePath := "test.csv"
			defer file.Remove()

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			routes := []r.Route{
				{
					Boarding:    "XYZ",
					Destination: "ABC",
					Cost:        1000,
				},
				{
					Boarding:    "XYZ",
					Destination: "EFG",
					Cost:        10,
				},
				{
					Boarding:    "XYZ",
					Destination: "ABC",
					Cost:        1001,
				},
				{
					Boarding:    "ab ",
					Destination: "ABC",
					Cost:        0,
				},
			}

			LoadRoutes(routes)

			cost, _ := database.GetRouteCost(routes[0].Boarding, routes[0].Destination)
			cost2, _ := database.GetRouteCost(routes[1].Boarding, routes[1].Destination)
			airport1 := database.GetAirport(routes[0].Boarding)
			airport2 := database.GetAirport(routes[0].Destination)
			airport3 := database.GetAirport(routes[1].Destination)
			routesFromCache := cache.GetAllRoutes()
			routesFromFile, _ := file.ReadFile()

			g.Assert(cost).Equal(1000)
			g.Assert(cost2).Equal(10)
			g.Assert(airport1).IsTrue()
			g.Assert(airport2).IsTrue()
			g.Assert(airport3).IsTrue()
			g.Assert(len(routesFromCache[routes[0].Boarding])).Equal(2)
			g.Assert(len(routesFromFile)).Equal(0)

			g.Assert(len(routesFromCache)).Equal(1)
			g.Assert(len(database.GetAllAirports())).Equal(3)
		})
	})

	g.Describe("Tests for GetBestRoute", func() {
		// Tests based on:
		//   GRU,BRC,10
		//   BRC,SCL,5
		//   GRU,CDG,75
		//   GRU,SCL,20
		//   GRU,ORL,56
		//   ORL,CDG,5
		//   SCL,ORL,20

		routes := []r.Route{
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
		g.It("should get the best route for every desired route", func() {
			filePath := "test.csv"
			defer file.Remove()

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			for _, route := range routes {
				err := AddNewRoute(route)
				g.Assert(err == nil).IsTrue()
			}

			best, _ := GetBestRoute("GRU", "BRC")
			best2, _ := GetBestRoute("GRU", "SCL")
			best3, _ := GetBestRoute("GRU", "ORL")
			best4, _ := GetBestRoute("GRU", "CDG")
			best5, _ := GetBestRoute("BRC", "SCL")
			best6, _ := GetBestRoute("SCL", "ORL")
			best7, _ := GetBestRoute("ORL", "CDG")
			best8, _ := GetBestRoute("BRC", "CDG")
			best9, _ := GetBestRoute("BRC", "ORL")
			best10, _ := GetBestRoute("SCL", "CDG")

			g.Assert(best.Route).Equal("GRU - BRC")
			g.Assert(best.Cost).Equal(10)

			g.Assert(best2.Route).Equal("GRU - BRC - SCL")
			g.Assert(best2.Cost).Equal(15)

			g.Assert(best3.Route).Equal("GRU - BRC - SCL - ORL")
			g.Assert(best3.Cost).Equal(35)

			g.Assert(best4.Route).Equal("GRU - BRC - SCL - ORL - CDG")
			g.Assert(best4.Cost).Equal(40)

			g.Assert(best5.Route).Equal("BRC - SCL")
			g.Assert(best5.Cost).Equal(5)

			g.Assert(best6.Route).Equal("SCL - ORL")
			g.Assert(best6.Cost).Equal(20)

			g.Assert(best7.Route).Equal("ORL - CDG")
			g.Assert(best7.Cost).Equal(5)

			g.Assert(best8.Route).Equal("BRC - SCL - ORL - CDG")
			g.Assert(best8.Cost).Equal(30)

			g.Assert(best9.Route).Equal("BRC - SCL - ORL")
			g.Assert(best9.Cost).Equal(25)

			g.Assert(best10.Route).Equal("SCL - ORL - CDG")
			g.Assert(best10.Cost).Equal(25)
		})

		g.It("should return InvalidAirportErr when an airport is not stored or has invalid format", func() {
			filePath := "test.csv"
			defer file.Remove()

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			for _, route := range routes {
				err := AddNewRoute(route)
				g.Assert(err == nil).IsTrue()
			}

			_, err := GetBestRoute("SCL", "XY")
			g.Assert(err).Equal(errors.NewInvalidAirportErr("invalid"))

			_, err = GetBestRoute("SCL", "XYZ")

			g.Assert(err).Equal(errors.NewInvalidAirportErr("not registered"))
		})

		g.It("should return BestRouteNotFoundErr when a route is not possible to be found", func() {
			filePath := "test.csv"
			defer file.Remove()

			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
			file.Reset(filePath)

			newRoutes := append(routes, r.Route{
				Boarding:    "XYZ",
				Destination: "ORL",
				Cost:        50,
			})

			for _, route := range newRoutes {
				err := AddNewRoute(route)
				g.Assert(err == nil).IsTrue()
			}

			_, err := GetBestRoute("SCL", "XYZ")

			g.Assert(err).Equal(errors.NewBestRouteNotFoundErr())
		})
	})
}
