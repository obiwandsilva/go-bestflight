# go-bestflight

[![Build Status](https://dev.azure.com/wandersonolivs/go-bestflight/_apis/build/status/obiwandsilva.go-bestflight?branchName=master)](https://dev.azure.com/wandersonolivs/go-bestflight/_build/latest?definitionId=2&branchName=master)

In order to be efficient when looking for the best routes for the registered connections, Bestflight use the [Dijkstra's Shortest Path algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm).

## Build

If you have `go >= 1.13`, you can build the application using the command `go mod download` and then `go build -o bestflight cmd/bestflight/main.go `.

Or you can run the pre built file bestflight in this repo with `./bestflight sourcefile.csv 5000`.

## Usage

go-bestflight has two interface that can be used: one `cli` to get the best routes via command line and an `HTTP API` for both get the best routes and register new routes.

An source file with pre-registered routes can be passed by argument. It must contain the following format:

```csv
GRU,BRC,10
BRC,SCL,5
GRU,CDG,75
GRU,SCL,20
GRU,ORL,56
ORL,CDG,5
SCL,ORL,20
```

When running the app, simply pass the file path as the first argument followed by the port the HTTP web server will use to run. Example:

    ./bestflight routes.csv 5000

An input will be asked:

    please enter the route:

Enter a desired boarding and destination in the format: `GRU-CDG`

An output will be given in the format: `best route: SCL - GRU - BRC > $25`

## API

The API has only two endpoint, one to register new routes and another to get the best route between two airports.

**Register new routes**

Method: *POST*

Endpoint: */route*

Contet-Type: *json*

Body Parameters:
 - *boarding:* string containing an airport with the format "GRU". Case insensitive.
 - *destination:* string containing an airport with the format "GRU". Case insensitive.
 - *cost:* integer with minum value of 0 and maximum of 1000000

Example:
```json
{
	"boarding": "SCL",
	"destination": "GRU",
	"cost": 15
}
```
Status Codes:
 - *201*: if successfully created
 - *200*: if the route already exists
 - *400*: malformed route

Response Body: same content sent.

**Get the best route between two airports**

Method: *GET*

Endpoint: */route*

Query Parameters:

 - *board*: string containing an airport with the format "GRU". Case insensitive.
 - *dest*: string containing an airport with the format "GRU". Case insensitive.

Example:
    
    /routes?board=SCL&dest=BRC

Status Codes:

 - *200*: if successfully found
 - *204*: searched, but not found
 - *400*: malformed route

Response body:
 - *route*: string containing the route in a readable way. Example: `SCL - GRU - BRC`
 - *cost*: integer with the value fot taking the route.

Example:
```json
{
    "route": "SCL - GRU - BRC",
    "cost": 25
}
```

## Docker

The application can also be executed in a container if you have `docker`. Follow the steps:

 - Create a file called `input.csv` into the project root.
 - Build the image: `docker build -t bestflight .`
 - Run a container with an interactive tty: `docker container run --name bestflight -it -p 5000:5000 bestflight sh`
 - Run the command: `./bestflight input.csv 5000`

## Application structure

The application structure is loosely based on the model [Ports and Adapters](https://dev.to/jofisaes/hexagonal-architecture-ports-and-adapters-1h4m). In real world usage, the best way of taking advantage of this strcture is by making have use of interfaces, but since this a tiny project, I didn't feel the for it.

The package structure is divided in three main components:

 - **application**: where different instances and interfaces of the application can be managed, such as this one with a cli and http interfaces.
        Also responsible for the appplication boostratp and configuration.
 - **domain**: contains all the entities, custom errors and services used for the business logic.
 - **resources**: any external resource used by the services in the domain, for example. Components like repositories
        as whole (file, dbs, etc), gateways (http clients) and others.

Although this project does not use databases or cache services, what is the ideal in a real world production environment, there are
three resources in this project that simulates then on memory. They are all under resources: cache, file anddatabase.

The cache resource simulates the usage of a Redis-like service where we could use as a routes-structure-ready-for-search, so we could
have an alternative option to the slow IO operations with the database.

The database is to simulate persist and more reliable data.

And finally the file resource to manage the input/source file used when running the application.
