package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

func draw(width, height int, grid [][]seatInfo) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetOutputMode(termbox.OutputNormal)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var color termbox.Attribute
			switch grid[y][x].state {
			case uncertain:
				color = termbox.ColorYellow
			case vacant:
				color = termbox.ColorGreen
			case floor:
				color = termbox.ColorBlue
			case occupied:
				color = termbox.ColorRed
			}
			termbox.SetCell(x, y, '▄', color, color)
		}
	}
	termbox.Flush()
}

type seatState int

const (
	floor     seatState = 0
	uncertain seatState = 1
	occupied  seatState = 2
	vacant    seatState = 3
)

type seatInfo struct {
	numNeighbors int
	state        seatState
}

type coord struct {
	x, y int
}

func visitNeighbors(width, height int, origin coord, job func(coord, coord)) {
	for y := -1; y <= 1; y++ {
		for x := -1; x <= 1; x++ {
			neighborY := origin.y + y
			neighborX := origin.x + x
			if 0 <= neighborX && neighborX < width && 0 <= neighborY && neighborY < height && !(y == 0 && x == 0) {
				job(origin, coord{neighborX, neighborY})
			}
		}
	}
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in11.txt", "Input file to use")
	flag.Parse()
	bytes, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(bytes), "\n") // parse input

	width, height := len(lines[0]), len(lines)
	gridMemory := make([]seatInfo, width*height)
	grid := make([][]seatInfo, height)
	for y := 0; y < height; y++ {
		grid[y] = gridMemory[y*width : (y+1)*width] // create 2D grid from linear memory
	}

	changes, newChanges := make([]coord, 0, width*height), make([]coord, 0, width*height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if code := lines[y][x]; code == 'L' { // Default is floor
				grid[y][x].state = uncertain

				// Count all neighbor seats for uncertain seats
				visitNeighbors(width, height, coord{x, y}, func(origin coord, n coord) {
					if neighborCode := lines[n.y][n.x]; neighborCode != '.' { // Detect neighbor seat
						grid[origin.y][origin.x].numNeighbors++
					}
				})

				if grid[y][x].numNeighbors < 4 { // Will definitely be occupied!
					grid[y][x].state = occupied
					changes = append(changes, coord{x, y})
				}
			}
		}
	}

	totalOccupied := len(changes) // Count initially occupied

	termbox.Init()
	// Only iterate over PERMANENT state changes. Keep going until nothing has changed anymore.
	// Update neighbors of previous permanent changes and record new permanent changes
	for len(changes) > 0 {
		for _, change := range changes {
			visitNeighbors(width, height, change, func(origin coord, n coord) {
				if grid[origin.y][origin.x].state == occupied && grid[n.y][n.x].state == uncertain {
					grid[n.y][n.x].state = vacant
					newChanges = append(newChanges, n)
				}
			})
		}

		for _, change := range changes {
			visitNeighbors(width, height, change, func(origin coord, n coord) {
				if grid[origin.y][origin.x].state == vacant && grid[n.y][n.x].state == uncertain {
					grid[n.y][n.x].numNeighbors--

					if grid[n.y][n.x].numNeighbors < 4 {
						grid[n.y][n.x].state = occupied
						newChanges = append(newChanges, n)
						totalOccupied++
					}
				}
			})
		}

		draw(width, height, grid)
		time.Sleep(50 * time.Millisecond)
		changes, newChanges = newChanges, changes[:0]
	}
	fmt.Printf("Total occupied: %d\n", totalOccupied)
}
