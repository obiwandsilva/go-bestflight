package routes

// Route ...
type Route struct {
	Boarding    string
	Destination string
	Cost        int
}

// BestRoute ...
type BestRoute struct {
	Route string
	Cost  int
}

// Connection ...
type Connection struct {
	Airport string
	Cost    int
}

// Routes ...
type Routes map[string][]Connection
