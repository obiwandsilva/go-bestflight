package validationservice

import (
	"testing"

	"github.com/franela/goblin"
)

func TestFile(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("tests for IsValidAirport", func() {
		g.It("should return true for valid airport formats and false for invalid ones", func() {
			g.Assert(IsValidAirport("ABC")).IsTrue()

			g.Assert(IsValidAirport("GR")).IsFalse()
			g.Assert(IsValidAirport("GRUU")).IsFalse()
			g.Assert(IsValidAirport("grU")).IsFalse()
			g.Assert(IsValidAirport("abc")).IsFalse()
			g.Assert(IsValidAirport("1BC")).IsFalse()
			g.Assert(IsValidAirport("@!#")).IsFalse()
			g.Assert(IsValidAirport(" ")).IsFalse()
			g.Assert(IsValidAirport("")).IsFalse()
		})
	})
}
