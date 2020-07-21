package cli

import (
	"bufio"
	"fmt"
	"go-bestflight/domain/services/routeservice"
	"log"
	"os"
	"runtime"
	"strings"
)

const delimiter = '\n'

func getInput() string {
	inputDeviceLocation := "/dev/stdin"

	if runtime.GOOS == "windows" {
		inputDeviceLocation = "CON"
	}

	file, err := os.Open(inputDeviceLocation)
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)
	input, err := reader.ReadString(delimiter)
	if err != nil {
		log.Fatal(err)
	}

	input = strings.Replace(input, "\n", "", -1)
	input = strings.Replace(input, "\r", "", -1) // for Windows
	input = strings.Trim(input, " ")

	return input
}

func getBoardingAndDestination(input string) (string, string) {
	components := strings.Split(input, "-")

	if len(components) != 2 {
		return "", ""
	}

	return components[0], components[1]
}

// StartAdvisor starts the agent that will be asking for desired routes by command line.
func StartAdvisor() {
	log.Println("starting Advisor...")

	for {
		fmt.Print("please enter the route: ")
		input := getInput()
		board, dest := getBoardingAndDestination(input)

		bestRoute, err := routeservice.GetBestRoute(board, dest)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("best route: %s > $%d\n", bestRoute.Route, bestRoute.Cost)
	}
}
