package database

// RouteNotFoundErr it's an error for when routes are not found when
// searched.
type RouteNotFoundErr struct {
	message string
}

func (dtf *RouteNotFoundErr) Error() string {
	return dtf.message
}

// NewRouteNotFoundErr ...
func NewRouteNotFoundErr() *RouteNotFoundErr {
	return &RouteNotFoundErr{
		message: "route not found",
	}
}
