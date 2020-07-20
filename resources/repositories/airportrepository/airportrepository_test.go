package airportrepository

import (
	"go-bestflight/resources/database"
	"testing"

	"github.com/franela/goblin"
)

func TestAirportRepository(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Tests for IsRegistered", func() {
		g.It("shoudl return false for a not stored airport and true for a stored one", func() {
			database.Connect()

			airport := "airportA"

			g.Assert(IsRegistered(airport)).IsFalse()

			database.StoreAirport(airport)

			g.Assert(IsRegistered(airport)).IsTrue()
		})
	})
}
