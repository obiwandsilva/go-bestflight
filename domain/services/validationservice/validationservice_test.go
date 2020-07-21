package validationservice

import (
	"go-bestflight/domain/entities/routes"
	"testing"

	"github.com/franela/goblin"
)

func TestFile(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Tests for isValidAirport", func() {
		g.It("should return true for valid airport formats and false for invalid ones", func() {
			g.Assert(isValidAirport("ABC")).IsTrue()

			g.Assert(isValidAirport("GR")).IsFalse()
			g.Assert(isValidAirport("GRUU")).IsFalse()
			g.Assert(isValidAirport("grU")).IsFalse()
			g.Assert(isValidAirport("abc")).IsFalse()
			g.Assert(isValidAirport("1BC")).IsFalse()
			g.Assert(isValidAirport("@!#")).IsFalse()
			g.Assert(isValidAirport(" ")).IsFalse()
			g.Assert(isValidAirport("")).IsFalse()
		})
	})

	g.Describe("Tests for isValidCost", func() {
		g.It("should return true for valid cost values and false for invalid ones", func() {
			g.Assert(isValidCost(1)).IsTrue()
			g.Assert(isValidCost(10)).IsTrue()
			g.Assert(isValidCost(100)).IsTrue()
			g.Assert(isValidCost(1000)).IsTrue()
			g.Assert(isValidCost(100000)).IsTrue()
			g.Assert(isValidCost(1000000)).IsTrue()

			g.Assert(isValidCost(0)).IsFalse()
			g.Assert(isValidCost(-1)).IsFalse()
			g.Assert(isValidCost(1000001)).IsFalse()
		})
	})

	g.Describe("Tests for IsValidRoute", func() {
		g.It("should return true for valid route formats and false for invalid ones", func() {
			g.Assert(
				IsValidRoute(routes.Route{
					Boarding:    "ABC",
					Destination: "ZYX",
					Cost:        3,
				}),
			).IsTrue()

			g.Assert(
				IsValidRoute(routes.Route{
					Boarding:    "",
					Destination: "ZYX",
					Cost:        3,
				}),
			).IsFalse()

			g.Assert(
				IsValidRoute(routes.Route{
					Boarding:    "abc",
					Destination: "ZYX",
					Cost:        3,
				}),
			).IsFalse()

			g.Assert(
				IsValidRoute(routes.Route{
					Boarding:    "ABCD",
					Destination: "ZYX",
					Cost:        3,
				}),
			).IsFalse()

			g.Assert(
				IsValidRoute(routes.Route{
					Boarding:    "ABC",
					Destination: "ZY",
					Cost:        3,
				}),
			).IsFalse()

			g.Assert(
				IsValidRoute(routes.Route{
					Boarding:    "ABC",
					Destination: "ZYX",
					Cost:        -1,
				}),
			).IsFalse()

			g.Assert(
				IsValidRoute(routes.Route{
					Boarding:    "ABC",
					Destination: "ZY",
					Cost:        1000001,
				}),
			).IsFalse()
		})
	})
}
