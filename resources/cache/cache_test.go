package cache

import (
	r "go-bestflight/domain/entities/routes"
	"testing"

	"github.com/franela/goblin"
)

func TestCache(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Tests for Connect", func() {
		g.It("should create only one instance on multiple calls", func() {
			Connect()

			instance.routes["GRU"] = []r.Connection{{Airport: "CDG", Cost: 5}}

			Connect()

			connections, ok := instance.routes["GRU"]

			g.Assert(ok).Equal(true)
			g.Assert(connections).Equal([]r.Connection{{Airport: "CDG", Cost: 5}})

			Truncate()
		})
	})

	g.Describe("Tests for AddRoute", func() {
		g.BeforeEach(func() {
			Connect()
		})

		g.AfterEach(func() {
			Truncate()
		})

		g.It("should successfully insert a new route to the cache", func() {
			AddRoute(r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			})

			for k, v := range instance.routes {
				g.Assert(k).Equal("GRU")
				g.Assert(v[0]).Equal(r.Connection{Airport: "CDG", Cost: 75})
			}
		})

		g.It("should insert multiple routes with same boarding and different destinations", func() {
			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := r.Route{
				Boarding:    "GRU",
				Destination: "ORL",
				Cost:        56,
			}

			AddRoute(route)
			AddRoute(route2)

			g.Assert(instance.routes["GRU"][0]).Equal(
				r.Connection{Airport: route.Destination, Cost: route.Cost},
			)
			g.Assert(instance.routes["GRU"][1]).Equal(
				r.Connection{Airport: route2.Destination, Cost: route2.Cost},
			)
		})

		g.It("should insert multiple routes with different boarding and destination", func() {
			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := r.Route{
				Boarding:    "SCL",
				Destination: "ORL",
				Cost:        20,
			}

			AddRoute(route)
			AddRoute(route2)

			g.Assert(instance.routes["GRU"][0]).Equal(
				r.Connection{Airport: route.Destination, Cost: route.Cost},
			)
			g.Assert(instance.routes["SCL"][0]).Equal(
				r.Connection{Airport: route2.Destination, Cost: route2.Cost},
			)
		})
	})

	g.Describe("Tests for GetAllRoutes", func() {
		g.It("should insert multiple routes from a list", func() {
			Connect()

			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := r.Route{
				Boarding:    "SCL",
				Destination: "ORL",
				Cost:        20,
			}
			route3 := r.Route{
				Boarding:    "BRC",
				Destination: "SCL",
				Cost:        5,
			}

			AddRoutes([]r.Route{route, route2, route3})

			g.Assert(instance.routes["GRU"][0]).Equal(
				r.Connection{Airport: route.Destination, Cost: route.Cost},
			)
			g.Assert(instance.routes["SCL"][0]).Equal(
				r.Connection{Airport: route2.Destination, Cost: route2.Cost},
			)
			g.Assert(instance.routes["BRC"][0]).Equal(
				r.Connection{Airport: route3.Destination, Cost: route3.Cost},
			)

			Truncate()
		})
	})

	g.Describe("Tests for GetAllRoutes", func() {
		g.It("should successfully return a copy of all current routes in cache", func() {
			Connect()

			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := r.Route{
				Boarding:    "SCL",
				Destination: "ORL",
				Cost:        20,
			}
			route3 := r.Route{
				Boarding:    "BRC",
				Destination: "SCL",
				Cost:        5,
			}

			AddRoutes([]r.Route{route, route2, route3})

			routes := GetAllRoutes()

			g.Assert(routes["GRU"][0]).Equal(
				r.Connection{Airport: route.Destination, Cost: route.Cost},
			)
			g.Assert(routes["SCL"][0]).Equal(
				r.Connection{Airport: route2.Destination, Cost: route2.Cost},
			)
			g.Assert(routes["BRC"][0]).Equal(
				r.Connection{Airport: route3.Destination, Cost: route3.Cost},
			)

			Truncate()
		})
	})
}
