package errors

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

// NewBestRouteNotFoundErr is aconstructor for BestRouteNotFoundErr.
func NewBestRouteNotFoundErr() *BestRouteNotFoundErr {
	return &BestRouteNotFoundErr{
		message: "best route not found",
	}
}
