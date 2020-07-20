package file

import (
	"bufio"
	"errors"
	"fmt"
	r "go-bestflight/domain/routes"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// RoutesFile ...
type RoutesFile struct {
	filePath string
	sync.RWMutex
}

var (
	instance RoutesFile
	once     sync.Once
)

// Sync ...
func Sync(filePath string) {
	once.Do(func() {
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0666)
		if err != nil {
			log.Fatal(err)
		}

		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}

		instance = RoutesFile{
			filePath: filePath,
		}
	})
}

// Only for tests.
func remove(filePath string) {
	os.Remove(filePath)
}

// Only for tests.
func truncate() {
	os.Truncate(instance.filePath, 0)
}

func isValidAirport(airport string) bool {
	pattern := `^[A-Z]{3}$`
	match, err := regexp.MatchString(pattern, airport)
	if err != nil {
		log.Println(err)
	}

	return match
}

func cleanLine(line string) string {
	return strings.Replace(line, "\n", "", -1)
}

func lineToRoute(line string, lineN int) (r.Route, error) {
	components := strings.Split(cleanLine(line), ",")

	if len(components) != 3 {
		log.Printf("invalid format at line: %d\n", lineN)
		return r.Route{}, errors.New("invalid line format")
	}

	boarding := strings.ToUpper(components[0])
	destination := strings.ToUpper(components[1])

	if !isValidAirport(boarding) || !isValidAirport(destination) {
		log.Printf("invalid airport format at line: %d\n", lineN)
		return r.Route{}, errors.New("invalid airport format")
	}

	cost, err := strconv.Atoi(components[2])
	if err != nil {
		log.Printf("error at line %d: %v\n", lineN, err)
		return r.Route{}, err
	}

	route := r.Route{
		Boarding:    boarding,
		Destination: destination,
		Cost:        cost,
	}

	return route, nil
}

func isEmpty(line string) bool {
	return len(line) == 1 && []byte(line)[0] == '\n'
}

// Write ...
func Write(route r.Route) error {
	instance.Lock()
	defer instance.Unlock()

	file, err := os.OpenFile(instance.filePath, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		log.Println(err)
		return err
	}

	strLine := fmt.Sprintf("%s,%s,%d\n", route.Boarding, route.Destination, route.Cost)

	_, err = file.WriteString(strLine)
	if err != nil {
		log.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// ReadFile ...
func ReadFile() ([]r.Route, error) {
	instance.RLock()
	defer instance.RUnlock()

	routes := []r.Route{}
	file, err := os.OpenFile(instance.filePath, os.O_RDONLY, 0444)
	if err != nil {
		log.Println(err)
		return routes, err
	}

	scan := bufio.NewScanner(file)
	lineNumber := 0

	for scan.Scan() {
		lineNumber++
		line := scan.Text()

		if isEmpty(line) {
			continue
		}

		route, err := lineToRoute(line, lineNumber)
		if err == nil {
			routes = append(routes, route)
		}
	}

	err = file.Close()
	if err != nil {
		log.Println(err)
		return routes, err
	}

	return routes, nil
}
