package errors

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
