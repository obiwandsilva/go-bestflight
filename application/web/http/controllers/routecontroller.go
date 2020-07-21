package controllers

import (
	"errors"
	r "go-bestflight/domain/entities/routes"
	e "go-bestflight/domain/errors"
	"go-bestflight/domain/services/routeservice"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddNewRoute(ctx *gin.Context) {
	var newRoute r.Route

	if err := ctx.ShouldBindJSON(&newRoute); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
	}

	addedRoute, err := routeservice.AddNewRoute(newRoute)
	if err != nil {
		if errors.Is(err, &e.InvalidRouteErr{}) {
			ctx.String(http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, &e.RouteAlreadyExistErr{}) {
			ctx.JSON(http.StatusCreated, newRoute)
		}

		ctx.String(http.StatusInternalServerError, "InternalServer Error")

		return
	}

	ctx.JSON(http.StatusCreated, addedRoute)
}

func BestRoute(ctx *gin.Context) {
	boarding := ctx.Query("board")
	destination := ctx.Query("dest")
	bestRoute, err := routeservice.GetBestRoute(boarding, destination)

	log.Println("HEREEEE", bestRoute, err)

	if err != nil {
		if errors.Is(err, &e.InvalidAirportErr{}) {
			ctx.String(http.StatusBadRequest, err.Error())
		}

		if errors.Is(err, &e.BestRouteNotFoundErr{}) {
			ctx.String(http.StatusBadRequest, err.Error())
		}

		return
	}

	ctx.JSON(http.StatusOK, bestRoute)
}
