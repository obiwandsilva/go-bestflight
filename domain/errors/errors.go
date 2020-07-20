package errors

import "fmt"

type InvalidRouteErr struct {
	message string
}

func (e *InvalidRouteErr) Error() string {
	return e.message
}

func NewInvalidRouteErr(cause string) *InvalidAirport {
	return &InvalidAirport{
		message: fmt.Sprintf("invalid: %s", cause),
	}
}

type InvalidAirport struct {
	message string
}

func (e *InvalidAirport) Error() string {
	return e.message
}

// NewInvalidAirportErr is a constructor for InvalidAirport. You can specify
// the cause as a parameter.
func NewInvalidAirportErr(cause string) *InvalidAirport {
	return &InvalidAirport{
		message: fmt.Sprintf("airport: %s", cause),
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
