package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

func draw(width, height int, grid [][]seatInformation) {
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

type seatSituation int

const (
	floor      seatSituation = 0
	uncertain  seatSituation = 1
	vacant     seatSituation = 2
	occupied   seatSituation = 4
)

type seatInformation struct {
	possibleNeighbors int
	state             seatSituation
	neighbors         []coord
}

type coord struct {
	x, y int
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in11.txt", "Input file to use")
	flag.Parse()
	bytes, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(bytes), "\n")

	width, height := len(lines[0]), len(lines)
	gridMemory := make([]seatInformation, width*height)
	grid := make([][]seatInformation, height)
	for y := 0; y < height; y++ {
		grid[y] = gridMemory[y*width : (y+1)*width]
	}

	updaters, nextUpdaters := make([]coord, width*height), make([]coord, width*height)
	numUpdaters, nextNumUpdaters := 0, 0
	totalOccupied := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			for ny := -1; ny < 2; ny++ {
				for nx := -1; nx < 2; nx++ {
					testY := y + ny
					testX := x + nx
					if 0 <= testY && testY < height && 0 <= testX && testX < width {
						seatCode := lines[testY][testX]
						var seat seatSituation
						switch seatCode {
						case '.':
							seat = floor
						case 'L':
							seat = uncertain
						}
						if ny == 0 && nx == 0 {
							grid[y][x].state = seat
						} else if seat == uncertain {
							grid[y][x].neighbors = append(grid[y][x].neighbors, coord{testX, testY})
							grid[y][x].possibleNeighbors++
						}
					}
				}
			}
			updaters[numUpdaters] = coord{x, y}
			numUpdaters++
		}
	}

	termbox.Init()

	for numUpdaters > 0 {
		nextNumUpdaters = 0
		for i := 0; i < numUpdaters; i++ {
			y, x := updaters[i].y, updaters[i].x
			if grid[y][x].state == occupied {
				for _, n := range grid[y][x].neighbors {
					if grid[n.y][n.x].state == uncertain {
						grid[n.y][n.x].state = vacant
						nextUpdaters[nextNumUpdaters] = n
						nextNumUpdaters++
					}
				}
			}
		}
		for i := 0; i < numUpdaters; i++ {
			y, x := updaters[i].y, updaters[i].x
			for _, n := range grid[y][x].neighbors {
				if grid[n.y][n.x].state == uncertain {
					if grid[y][x].state == vacant {
						grid[n.y][n.x].possibleNeighbors--
					}
					if grid[n.y][n.x].possibleNeighbors < 4 {
						grid[n.y][n.x].state = occupied
						nextUpdaters[nextNumUpdaters] = n
						nextNumUpdaters++
						totalOccupied++
					}
				}
			}
		}
		draw(width, height, grid)
		time.Sleep(50 * time.Millisecond)
		numUpdaters, nextNumUpdaters = nextNumUpdaters, numUpdaters
		updaters, nextUpdaters = nextUpdaters, updaters
	}
	fmt.Printf("Total occupied: %d\n", totalOccupied)
}
