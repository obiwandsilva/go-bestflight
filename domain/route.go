package domain

// Route ...
type Route struct {
	Boarding    string
	Destination string
	Cost        int
}

// Destination ...
type Destination struct {
	Airport string
	Cost    int
}

// Routes ...
type Routes struct {
	routes map[string][]Destination
}
