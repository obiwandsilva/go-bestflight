package cli

import (
	"bufio"
	"fmt"
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

// StartAdvisor starts the agent that will be asking for desired routes by command line.
func StartAdvisor() {
	for {
		fmt.Print("please enter the route: ")
		getInput()
	}
}
