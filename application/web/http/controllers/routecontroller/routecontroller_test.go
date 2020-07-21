package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	r "go-bestflight/domain/entities/routes"
	"go-bestflight/domain/errors"
	"go-bestflight/resources/cache"
	"go-bestflight/resources/database"
	"go-bestflight/resources/file"
	"go-bestflight/resources/repositories/routerepository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func TestController(t *testing.T) {
	g := goblin.Goblin(t)

	gin.SetMode(gin.TestMode)

	g.Describe("Tests for AddNewRoute", func() {
		g.It("should add a new route and return status code 201 and a json with the route info", func() {
			filePath := "test.csv"
			defer file.Remove()

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			jsonBytes, _ := json.Marshal(route)
			req, _ := http.NewRequest("POST", "localhost:3000/route", bytes.NewReader(jsonBytes))
			resWriter := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resWriter)
			ctx.Request = req

			g.Assert(resWriter.Body.String()).Equal("")
			g.Assert(routerepository.RouteExists(route.Boarding, route.Destination)).IsFalse()

			AddNewRoute(ctx)

			g.Assert(routerepository.RouteExists(route.Boarding, route.Destination)).IsTrue()
			g.Assert(resWriter.Code).Equal(201)
			g.Assert(resWriter.Body.String()).Equal(string(jsonBytes))
		})

		g.It("should return status code 400 for invalid data", func() {
			filePath := "test.csv"
			defer file.Remove()

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			invalidRoute := []int{1, 2, 3}
			jsonBytes, _ := json.Marshal(invalidRoute)
			req, _ := http.NewRequest("POST", "localhost:3000/route", bytes.NewReader(jsonBytes))
			resWriter := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resWriter)
			ctx.Request = req

			g.Assert(resWriter.Body.String()).Equal("")

			AddNewRoute(ctx)

			g.Assert(resWriter.Code).Equal(400)
			g.Assert(resWriter.Body.String()).Equal("Bad Request")
		})

		g.It("should return status code 400 for malformed route", func() {
			filePath := "test.csv"
			defer file.Remove()

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			malformedRoute := r.Route{
				Boarding:    "",
				Destination: "CDG",
				Cost:        0,
			}
			jsonBytes, _ := json.Marshal(malformedRoute)
			req, _ := http.NewRequest("POST", "localhost:3000/route", bytes.NewReader(jsonBytes))
			resWriter := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resWriter)
			ctx.Request = req

			g.Assert(resWriter.Body.String()).Equal("")

			AddNewRoute(ctx)

			g.Assert(resWriter.Code).Equal(400)
			g.Assert(resWriter.Body.String()).Equal(errors.NewInvalidRouteErr().Error())
		})

		g.It("should return status code 200 for a already created route", func() {
			filePath := "test.csv"
			defer file.Remove()

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        200,
			}
			jsonBytes, _ := json.Marshal(route)
			req, _ := http.NewRequest("POST", "localhost:3000/route", bytes.NewReader(jsonBytes))
			resWriter := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resWriter)
			ctx.Request = req

			g.Assert(resWriter.Body.String()).Equal("")

			AddNewRoute(ctx)

			g.Assert(resWriter.Code).Equal(201)

			jsonBytes2, _ := json.Marshal(route)
			req2, _ := http.NewRequest("POST", "localhost:3000/route", bytes.NewReader(jsonBytes2))
			resWriter2 := httptest.NewRecorder()
			ctx2, _ := gin.CreateTestContext(resWriter2)
			ctx2.Request = req2

			g.Assert(resWriter2.Body.String()).Equal("")

			AddNewRoute(ctx2)

			g.Assert(resWriter2.Code).Equal(200)
			g.Assert(resWriter2.Body.String()).Equal(string(jsonBytes2))
		})
	})

	g.Describe("Tests for BestRoute", func() {
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
				req, _ := http.NewRequest("POST", "localhost:3000/route", bytes.NewReader(jsonBytes))
				resWriter := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(resWriter)
				ctx.Request = req

				AddNewRoute(ctx)
			}
		}

		g.It("should return status code 200 a json with best route info", func() {
			filePath := "test.csv"
			defer file.Remove()

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			addRoutes()

			board := "gru"
			dest := "cdg"
			req, _ := http.NewRequest("POST", fmt.Sprintf("localhost:3000/route?board=%s&dest=%s", board, dest), bytes.NewReader([]byte{}))
			resWriter := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resWriter)
			ctx.Request = req

			BestRoute(ctx)

			expectedBestRoute := r.BestRoute{
				Route: "GRU - BRC - SCL - ORL - CDG",
				Cost:  40,
			}
			jsonData, _ := json.Marshal(expectedBestRoute)

			g.Assert(resWriter.Code).Equal(200)
			g.Assert(resWriter.Body.String()).Equal(string(jsonData))
		})

		g.It("should return status code 400 for a not stored airport", func() {
			filePath := "test.csv"
			defer file.Remove()

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			addRoutes()

			notStoredBoard := "ACB"
			dest := "cdg"
			req, _ := http.NewRequest("POST", fmt.Sprintf("localhost:3000/route?board=%s&dest=%s", notStoredBoard, dest), bytes.NewReader([]byte{}))
			resWriter := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resWriter)
			ctx.Request = req

			BestRoute(ctx)

			g.Assert(resWriter.Code).Equal(400)
			g.Assert(resWriter.Body.String()).Equal(errors.NewInvalidAirportErr("not registered").Error())
		})

		g.It("should return status code 400 for a malformed airport", func() {
			filePath := "test.csv"
			defer file.Remove()

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			addRoutes()

			malformedBoard := "A"
			dest := "cdg"
			req, _ := http.NewRequest("POST", fmt.Sprintf("localhost:3000/route?board=%s&dest=%s", malformedBoard, dest), bytes.NewReader([]byte{}))
			resWriter := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resWriter)
			ctx.Request = req

			BestRoute(ctx)

			g.Assert(resWriter.Code).Equal(400)
			g.Assert(resWriter.Body.String()).Equal(errors.NewInvalidAirportErr("malformed").Error())
		})

		g.It("should return status code 204 when the best route is not found", func() {
			filePath := "test.csv"
			defer file.Remove()

			file.Reset(filePath)
			database.Connect()
			database.Truncate()
			cache.Connect()
			cache.Truncate()

			addRoutes()

			board := "SCL"
			dest := "gru"
			req, _ := http.NewRequest("POST", fmt.Sprintf("localhost:3000/route?board=%s&dest=%s", board, dest), bytes.NewReader([]byte{}))
			resWriter := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(resWriter)
			ctx.Request = req

			BestRoute(ctx)

			g.Assert(resWriter.Code).Equal(204)
		})
	})
}
