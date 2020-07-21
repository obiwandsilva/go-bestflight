package routes

// Route ...
type Route struct {
	Boarding    string `json:"boarding"`
	Destination string `json:"destination"`
	Cost        int    `json:"cost"`
}

// BestRoute ...
type BestRoute struct {
	Route string `json:"route"`
	Cost  int    `json:"cost"`
}

// Connection ...
type Connection struct {
	Airport string
	Cost    int
}

// Routes ...
type Routes map[string][]Connection
