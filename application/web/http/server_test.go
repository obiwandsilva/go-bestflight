package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	r "go-bestflight/domain/entities/routes"
	"go-bestflight/resources/cache"
	"go-bestflight/resources/database"
	"go-bestflight/resources/file"
	"go-bestflight/resources/repositories/routerepository"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func TestComponents(t *testing.T) {
	g := goblin.Goblin(t)

	Start("3000", gin.TestMode, nil)
	time.Sleep(2 * time.Second)

	g.Describe("Tests for the adding of new routes", func() {
		g.BeforeEach(func() {
			filePath := "test.csv"

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()
		})

		g.AfterEach(func() {
			file.Remove()
		})

		g.It("should return status code 201 when adding a new and valid route", func() {
			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			jsonBytes, _ := json.Marshal(route)

			resp, _ := http.Post("http://localhost:3000/routes", "application/json", bytes.NewReader(jsonBytes))
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			g.Assert(resp.StatusCode).Equal(201)
			g.Assert(jsonBytes).Equal(body)

			g.Assert(routerepository.RouteExists(route.Boarding, route.Destination)).IsTrue()
		})

		g.It("should return status code 200 when adding existing routes", func() {
			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			jsonBytes, _ := json.Marshal(route)

			http.Post("http://localhost:3000/routes", "application/json", bytes.NewReader(jsonBytes))

			resp, _ := http.Post("http://localhost:3000/routes", "application/json", bytes.NewReader(jsonBytes))

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			g.Assert(resp.StatusCode).Equal(200)
			g.Assert(jsonBytes).Equal(body)
		})

		g.It("should return status code 400 for a malformed route", func() {
			route := r.Route{
				Boarding:    "GRU",
				Destination: "DG",
				Cost:        0,
			}
			jsonBytes, _ := json.Marshal(route)

			resp, _ := http.Post("http://localhost:3000/routes", "application/json", bytes.NewReader(jsonBytes))

			g.Assert(resp.StatusCode).Equal(400)

			g.Assert(routerepository.RouteExists(route.Boarding, route.Destination)).IsFalse()
		})
	})

	g.Describe("Tests for getting best route", func() {
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

		addRoutes := func() {
			for _, route := range routes {
				jsonBytes, _ := json.Marshal(route)
				http.Post("http://localhost:3000/routes", "application/json", bytes.NewReader(jsonBytes))
			}
		}

		g.BeforeEach(func() {
			filePath := "test.csv"

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			addRoutes()
		})

		g.AfterEach(func() {
			file.Remove()
		})

		g.It("should return status code 200 and best route info json when successfully find a the best route", func() {
			boarding := "gru"
			destination := "CDG"
			expectedBestRoute := r.BestRoute{
				Route: "GRU - BRC - SCL - ORL - CDG",
				Cost:  40,
			}
			jsonData, _ := json.Marshal(expectedBestRoute)

			resp, _ := http.Get(fmt.Sprintf("http://localhost:3000/routes?board=%s&dest=%s", boarding, destination))
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)

			g.Assert(resp.StatusCode).Equal(200)
			g.Assert(jsonData).Equal(body)
		})

		g.It("should return status code 400 for malformed airport", func() {
			malformedBoarding := " R1"
			destination := "CDG"

			resp, _ := http.Get(fmt.Sprintf("http://localhost:3000/routes?board=%s&dest=%s", malformedBoarding, destination))

			g.Assert(resp.StatusCode).Equal(400)
		})

		g.It("should return status code 204 when the best route could not be found", func() {
			boarding := "SCL"
			destination := "GRU"

			resp, _ := http.Get(fmt.Sprintf("http://localhost:3000/routes?board=%s&dest=%s", boarding, destination))

			g.Assert(resp.StatusCode).Equal(204)
		})
	})
}
