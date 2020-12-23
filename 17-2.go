package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"time"
)

type int4 struct {
	x, y, z, w int
}

type Cube struct {
	location         int4
	on               bool
	nextOn           bool
	neighbors        [81]*Cube
	presentNeighbors int
}

var sceneMin = int4{math.MaxInt32, math.MaxInt32, math.MaxInt32, math.MaxInt32}
var sceneMax = int4{-math.MaxInt32, -math.MaxInt32, -math.MaxInt32, -math.MaxInt32}

func (A int4) Add(B int4) int4 {
	return int4{A.x + B.x, A.y + B.y, A.z + B.z, A.w + B.w}
}

func (A int4) Sub(B int4) int4 {
	return int4{A.x - B.x, A.y - B.y, A.z - B.z, A.w - B.w}
}

func max(a int4, b int4) int4 {
	res := b
	if a.x > b.x {
		res.x = a.x
	}
	if a.y > b.y {
		res.y = a.y
	}
	if a.z > b.z {
		res.z = a.z
	}
	if a.w > b.w {
		res.w = a.w
	}
	return res
}

func min(a int4, b int4) int4 {
	res := b
	if a.x < b.x {
		res.x = a.x
	}
	if a.y < b.y {
		res.y = a.y
	}
	if a.z < b.z {
		res.z = a.z
	}
	if a.w < b.w {
		res.w = a.w
	}
	return res
}

func (c *Cube) setNeighbor(coord int4, neighbor *Cube) {
	loc := (coord.w+1)*27 + (coord.z+1)*9 + (coord.y+1)*3 + coord.x + 1
	c.neighbors[loc] = neighbor
}

func (c *Cube) getNeighbor(coord int4) *Cube {
	loc := (coord.w+1)*27 + (coord.z+1)*9 + (coord.y+1)*3 + coord.x + 1
	return c.neighbors[loc]
}

var allCubes []*Cube
var incomplete = make(map[int4]*Cube)

func addMissingNeighborsOff(cube *Cube) {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				for w := -1; w <= 1; w++ {
					if x == 0 && y == 0 && z == 0 && w == 0 {
						continue
					}
					offset := int4{x, y, z, w}
					if cube.getNeighbor(offset) == nil {
						newCube := Cube{location: cube.location.Add(offset)}
						addCube(&newCube)
					}
				}
			}
		}
	}
}

func addCube(cube *Cube) {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				for w := -1; w <= 1; w++ {
					if x == 0 && y == 0 && z == 0 && w == 0 {
						continue
					}
					offset := int4{x, y, z, w}
					negOffset := int4{-x, -y, -z, -w}
					neighborCoord := cube.location.Add(offset)

					if neighbor, ok := incomplete[neighborCoord]; ok {
						cube.setNeighbor(offset, neighbor)
						cube.presentNeighbors++

						neighbor.setNeighbor(negOffset, cube)
						neighbor.presentNeighbors++

						if neighbor.presentNeighbors == 80 {
							delete(incomplete, neighborCoord)
						}
					}
				}
			}
		}
	}
	if cube.presentNeighbors < 80 {
		incomplete[cube.location] = cube
	}
	sceneMin = min(cube.location, sceneMin)
	sceneMax = max(cube.location, sceneMax)
	allCubes = append(allCubes, cube)
}

func main() {
	var fileName string
	var version2 bool
	flag.StringVar(&fileName, "file", "data/in17.txt", "Input file to use")
	flag.BoolVar(&version2, "v2", false, "Use task2 version")
	flag.Parse()

	content, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(content), "\n")

	before := time.Now()
	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				cube := Cube{on: true, location: int4{x, y, 0, 0}}
				addCube(&cube)
			}
		}
	}

	for _, cube := range allCubes {
		addMissingNeighborsOff(cube)
	}

	for i := 0; i < 6; i++ {
		for c := 0; c < len(allCubes); c++ {
			neighborsOn := 0
			for _, neighbor := range allCubes[c].neighbors {
				if neighbor != nil && neighbor.on {
					neighborsOn++
				}
			}
			if (neighborsOn == 3) || (allCubes[c].on && neighborsOn == 2) {
				allCubes[c].nextOn = true
				if _, ok := incomplete[allCubes[c].location]; ok {
					addMissingNeighborsOff(allCubes[c])
				}
			} else {
				allCubes[c].nextOn = false
			}
		}
		for _, cube := range allCubes {
			cube.on = cube.nextOn
		}
	}

	// Output

	// sceneDim := sceneMax.Sub(sceneMin).Add(int4{1, 1, 1, 1})
	// grid := make([]rune, sceneDim.x*sceneDim.y*sceneDim.z*sceneDim.w)
	// for w := 0; w < sceneDim.w; w++ {
	// 	for z := 0; z < sceneDim.z; z++ {
	// 		for y := 0; y < sceneDim.y; y++ {
	// 			for x := 0; x < sceneDim.x; x++ {
	// 				pos := w*(sceneDim.x*sceneDim.y*sceneDim.z) + z*(sceneDim.x*sceneDim.y) + y*sceneDim.x + x
	// 				grid[pos] = '.'
	// 			}
	// 		}
	// 	}
	// }

	// for _, cube := range allCubes {
	// 	if cube.on {
	// 		off := cube.location.Sub(sceneMin)
	// 		pos := off.w*(sceneDim.x*sceneDim.y*sceneDim.z) + off.z*(sceneDim.x*sceneDim.y) + off.y*sceneDim.x + off.x
	// 		grid[pos] = '#'
	// 	}
	// }

	// for w := 1; w < sceneDim.w-1; w++ {
	// 	for z := 1; z < sceneDim.z-1; z++ {
	// 		for y := 1; y < sceneDim.y-1; y++ {
	// 			for x := 1; x < sceneDim.x-1; x++ {
	// 				pos := w*(sceneDim.x*sceneDim.y*sceneDim.z) + z*(sceneDim.x*sceneDim.y) + y*sceneDim.x + x
	// 				fmt.Printf("%c", grid[pos])
	// 			}
	// 			fmt.Println("")
	// 		}
	// 		fmt.Println("")
	// 	}
	// }

	onSum := 0
	for _, c := range allCubes {
		if c.on {
			onSum++
		}
	}

	fmt.Println(time.Since(before).Milliseconds(), " ms")

	fmt.Println(onSum)
}
