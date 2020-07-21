package errors

import "fmt"

// InvalidRouteErr define errors with route format.
type InvalidRouteErr struct {
	message string
}

func (e *InvalidRouteErr) Error() string {
	return e.message
}

func NewInvalidRouteErr() *InvalidRouteErr {
	return &InvalidRouteErr{
		message: "invalid route format",
	}
}

type RouteAlreadyExistErr struct {
	message string
}

func (e *RouteAlreadyExistErr) Error() string {
	return e.message
}

func NewRouteAlreadyExistErr() *RouteAlreadyExistErr {
	return &RouteAlreadyExistErr{
		message: "route already created",
	}
}

type InvalidAirportErr struct {
	message string
}

func (e *InvalidAirportErr) Error() string {
	return e.message
}

func NewInvalidAirportErr(cause string) *InvalidAirportErr {
	return &InvalidAirportErr{
		message: fmt.Sprintf("invalid airport: %s", cause),
	}
}

// RouteNotFoundErr is an error for when routes are not found when searched.
type RouteNotFoundErr struct {
	message string
}

func (e *RouteNotFoundErr) Error() string {
	return e.message
}

// NewRouteNotFoundErr is aconstructor for RouteNotFoundErr.
func NewRouteNotFoundErr() *RouteNotFoundErr {
	return &RouteNotFoundErr{
		message: "route not found",
	}
}

// BestRouteNotFoundErr represents a searched but not found best route.
type BestRouteNotFoundErr struct {
	message string
}

func (e *BestRouteNotFoundErr) Error() string {
	return e.message
}

// NewBestRouteNotFoundErr is a constructor for BestRouteNotFoundErr.
func NewBestRouteNotFoundErr() *BestRouteNotFoundErr {
	return &BestRouteNotFoundErr{
		message: "best route not found",
	}
}
