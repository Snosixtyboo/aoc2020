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
	willVacate seatSituation = 8
	mayOccupy  seatSituation = 16
)

type seatInformation struct {
	possibleNeighbors int
	occupiedNeighbors int
	state             seatSituation
	updateState       seatSituation
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

	before := time.Now()

	updaters := make([]coord, width*height)
	numUpdaters := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			seatCode := lines[y][x]
			var seat seatSituation
			switch seatCode {
			case '.':
				seat = floor
			case 'L':
				seat = uncertain
			}
			grid[y][x].state = seat
			updaters[numUpdaters] = coord{x, y}
			numUpdaters++
		}
	}

	for y := 0; y < height; y++ {
		prev := coord{-1, -1}
		for x := 0; x < width; x++ {
			if grid[y][x].state == uncertain {
				if prev.x != -1 {
					grid[y][x].neighbors = append(grid[y][x].neighbors, prev)
					grid[y][x].possibleNeighbors++
					grid[prev.y][prev.x].neighbors = append(grid[prev.y][prev.x].neighbors, coord{x, y})
					grid[prev.y][prev.x].possibleNeighbors++
				}
				prev = coord{x, y}
			}
		}
	}

	for x := 0; x < width; x++ {
		prev := coord{-1, -1}
		for y := 0; y < height; y++ {
			if grid[y][x].state == uncertain {
				if prev.x != -1 {
					grid[y][x].neighbors = append(grid[y][x].neighbors, prev)
					grid[y][x].possibleNeighbors++
					grid[prev.y][prev.x].neighbors = append(grid[prev.y][prev.x].neighbors, coord{x, y})
					grid[prev.y][prev.x].possibleNeighbors++
				}
				prev = coord{x, y}
			}
		}
	}

	for a := 0; a < height+width-1; a++ {
		starty := a
		startx := 0
		if starty >= height {
			starty = 0
			startx = a - height + 1
		}
		prev := coord{-1, -1}
		for x, y := startx, starty; x < width && y < height; {
			if grid[y][x].state == uncertain {
				if prev.x != -1 {
					grid[y][x].neighbors = append(grid[y][x].neighbors, prev)
					grid[y][x].possibleNeighbors++
					grid[prev.y][prev.x].neighbors = append(grid[prev.y][prev.x].neighbors, coord{x, y})
					grid[prev.y][prev.x].possibleNeighbors++
				}
				prev = coord{x, y}
			}
			x++
			y++
		}
	}

	for a := 0; a < height+width-1; a++ {
		starty := a
		startx := 0
		if starty >= height {
			starty = height - 1
			startx = a - height + 1
		}
		prev := coord{-1, -1}
		for x, y := startx, starty; x < width && y >= 0; {
			if grid[y][x].state == uncertain {
				if prev.x != -1 {
					grid[y][x].neighbors = append(grid[y][x].neighbors, prev)
					grid[y][x].possibleNeighbors++
					grid[prev.y][prev.x].neighbors = append(grid[prev.y][prev.x].neighbors, coord{x, y})
					grid[prev.y][prev.x].possibleNeighbors++
				}
				prev = coord{x, y}
			}
			x++
			y--
		}
	}

	termbox.Init()
	totalOccupied := 0
	nextUpdaters := make([]coord, width*height)
	nextNumUpdaters := 0
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
					if grid[n.y][n.x].possibleNeighbors < 5 {
						grid[n.y][n.x].state = occupied
						nextUpdaters[nextNumUpdaters] = n
						nextNumUpdaters++
						totalOccupied++
					}
				}
			}
		}
		//draw(width, height, grid)
		//time.Sleep(50 * time.Millisecond)
		numUpdaters, nextNumUpdaters = nextNumUpdaters, numUpdaters
		updaters, nextUpdaters = nextUpdaters, updaters
	}

	fmt.Printf("Total occupied: %d\n", totalOccupied)
	fmt.Println(time.Since(before))
}
