package controllers

import (
	r "go-bestflight/domain/entities/routes"
	"go-bestflight/domain/errors"
	"go-bestflight/domain/services/routeservice"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddNewRoute is a handler for API route GET /route.
func AddNewRoute(ctx *gin.Context) {
	var newRoute r.Route

	if err := ctx.ShouldBindJSON(&newRoute); err != nil {
		ctx.String(http.StatusBadRequest, "Bad Request")

		return
	}

	addedRoute, err := routeservice.AddNewRoute(newRoute)
	if err != nil {
		if e, ok := err.(*errors.InvalidRouteErr); ok {
			ctx.String(http.StatusBadRequest, e.Error())
			return
		}

		if _, ok := err.(*errors.RouteAlreadyExistErr); ok {
			ctx.JSON(http.StatusOK, newRoute)
			return
		}

		log.Printf("unkown error when adding new route: %v", err)

		ctx.String(http.StatusInternalServerError, "Internal Server Error")

		return
	}

	ctx.JSON(http.StatusCreated, addedRoute)
}

// BestRoute is a handler for API route POST /route.
func BestRoute(ctx *gin.Context) {
	boarding := ctx.Query("board")
	destination := ctx.Query("dest")
	bestRoute, err := routeservice.GetBestRoute(boarding, destination)

	if err != nil {

		log.Println("####", err)

		if e, ok := err.(*errors.InvalidAirportErr); ok {
			ctx.String(http.StatusBadRequest, e.Error())
			return
		}

		if _, ok := err.(*errors.BestRouteNotFoundErr); ok {
			ctx.String(http.StatusNoContent, "")
			return
		}

		log.Printf("unkown error when getting best route: %v", err)

		ctx.String(http.StatusInternalServerError, "Internal Server Error")

		return
	}

	ctx.JSON(http.StatusOK, bestRoute)
}
