package file

import (
	"bufio"
	r "go-bestflight/domain/entities/routes"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/franela/goblin"
)

func TestFile(t *testing.T) {
	g := goblin.Goblin(t)

	filePath := "test.csv"

	g.Describe("Tests for ", func() {
		g.It("should create one instance on multiple calls and create file if it does not exist", func() {
			remove(filePath)

			_, err := os.Stat(filePath)
			g.Assert(os.IsNotExist(err)).IsTrue()

			Sync(filePath)
			Sync("anotherFile.csv")
			Sync("thirdFile.csv")

			g.Assert(instance.filePath).Equal(filePath)

			fileInfo, _ := os.Stat(filePath)

			g.Assert(fileInfo.Name()).Equal(filePath)
		})
	})

	g.Describe("Tests for lineToRoute", func() {
		g.It("should convert a valid line to a Route", func() {
			route, _ := lineToRoute("GRU,CDG,75", 1)

			g.Assert(route.Boarding).Equal("GRU")
			g.Assert(route.Destination).Equal("CDG")
			g.Assert(route.Cost).Equal(75)
		})

		g.It("should return error for invalid airport format", func() {
			_, err := lineToRoute("GR,CDG,75", 1)
			g.Assert(err != nil).IsTrue()

			_, err = lineToRoute("GRU,CD,75", 2)
			g.Assert(err != nil).IsTrue()

			_, err = lineToRoute("GRU,1DG,75", 3)
			g.Assert(err != nil).IsTrue()

			_, err = lineToRoute("GR U,CDG,75", 4)
			g.Assert(err != nil).IsTrue()
		})

		g.It("should return error if cost is not a valid integer", func() {
			_, err := lineToRoute("GRU,CDG,75.", 1)
			g.Assert(err != nil).IsTrue()

			_, err = lineToRoute("GRU,CDG,75.6", 1)
			g.Assert(err != nil).IsTrue()

			_, err = lineToRoute("GRU,CDG,", 1)
			g.Assert(err != nil).IsTrue()
		})
	})

	g.Describe("Tests for Write", func() {
		g.It("should successfully write to file", func() {
			truncate()

			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}

			Write(route)

			file, _ := os.Open(filePath)
			reader := bufio.NewReader(file)
			line, _ := reader.ReadString('\n')
			routeFromFile, _ := lineToRoute(line, 1)

			g.Assert(routeFromFile).Equal(route)
		})

		g.It("should successfully write to file concurrently", func() {
			truncate()

			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := r.Route{
				Boarding:    "GRU",
				Destination: "BRC",
				Cost:        10,
			}
			route3 := r.Route{
				Boarding:    "GRU",
				Destination: "SCL",
				Cost:        20,
			}
			route4 := r.Route{
				Boarding:    "GRU",
				Destination: "ORL",
				Cost:        56,
			}
			route5 := r.Route{
				Boarding:    "BRC",
				Destination: "SCL",
				Cost:        5,
			}

			var wg sync.WaitGroup
			wg.Add(5)

			go func() {
				Write(route)
				wg.Done()
			}()
			go func() {
				Write(route2)
				wg.Done()
			}()
			go func() {
				Write(route3)
				wg.Done()
			}()
			go func() {
				Write(route4)
				wg.Done()
			}()
			go func() {
				Write(route5)
				wg.Done()
			}()

			wg.Wait()

			file, err := os.Open(filePath)
			if err != nil {
				log.Fatal(err)
			}

			reader := bufio.NewReader(file)
			counter := 0

			for line, err := reader.ReadString('\n'); err == nil; line, err = reader.ReadString('\n') {
				_, errLine := lineToRoute(line, counter+1)
				if errLine == nil {
					counter++
				}
			}

			g.Assert(counter).Equal(5)
		})
	})

	g.Describe("Tests for ReadFile", func() {
		g.It("should return a list with all the written routes in the file", func() {
			truncate()

			route := r.Route{
				Boarding:    "GRU",
				Destination: "CDG",
				Cost:        75,
			}
			route2 := r.Route{
				Boarding:    "GRU",
				Destination: "BRC",
				Cost:        10,
			}
			route3 := r.Route{
				Boarding:    "GRU",
				Destination: "SCL",
				Cost:        20,
			}
			route4 := r.Route{
				Boarding:    "GRU",
				Destination: "ORL",
				Cost:        56,
			}
			route5 := r.Route{
				Boarding:    "BRC",
				Destination: "SCL",
				Cost:        5,
			}

			Write(route)
			Write(route2)
			Write(route3)
			Write(route4)
			Write(route5)

			routes, _ := ReadFile()

			g.Assert(routes[0]).Equal(route)
			g.Assert(routes[1]).Equal(route2)
			g.Assert(routes[2]).Equal(route3)
			g.Assert(routes[3]).Equal(route4)
			g.Assert(routes[4]).Equal(route5)
		})
	})

	remove(filePath)
}
