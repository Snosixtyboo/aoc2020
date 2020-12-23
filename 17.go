package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type int3 struct {
	x, y, z int
}

type Cube struct {
	location         int3
	on               bool
	nextOn           bool
	neighbors        [27]*Cube
	presentNeighbors int
}

var sceneMin = int3{math.MaxInt32, math.MaxInt32, math.MaxInt32}
var sceneMax = int3{-math.MaxInt32, -math.MaxInt32, -math.MaxInt32}

func (A int3) Add(B int3) int3 {
	return int3{A.x + B.x, A.y + B.y, A.z + B.z}
}

func (A int3) Sub(B int3) int3 {
	return int3{A.x - B.x, A.y - B.y, A.z - B.z}
}

func max(a int3, b int3) int3 {
	var res int3 = b
	if a.x > b.x {
		res.x = a.x
	}
	if a.y > b.y {
		res.y = a.y
	}
	if a.z > b.z {
		res.z = a.z
	}
	return res
}

func min(a int3, b int3) int3 {
	var res int3 = b
	if a.x < b.x {
		res.x = a.x
	}
	if a.y < b.y {
		res.y = a.y
	}
	if a.z < b.z {
		res.z = a.z
	}
	return res
}

func (c *Cube) setNeighbor(coord int3, neighbor *Cube) {
	loc := (coord.z+1)*9 + (coord.y+1)*3 + coord.x + 1
	c.neighbors[loc] = neighbor
}

func (c *Cube) getNeighbor(coord int3) *Cube {
	loc := (coord.z+1)*9 + (coord.y+1)*3 + coord.x + 1
	return c.neighbors[loc]
}

var allCubes []*Cube
var incomplete = make(map[int3]*Cube)

func addMissingNeighborsOff(cube *Cube) {
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				offset := int3{x, y, z}
				if cube.getNeighbor(offset) == nil {
					//fmt.Println("Going ", cube.location.x, cube.location.y, cube.location.z)
					//fmt.Println("Going ", offset.x, offset.y, offset.z)
					newCube := Cube{location: cube.location.Add(offset)}
					addCube(&newCube)
					//fmt.Println("\n")
				}
			}
		}
	}
}

func addCube(cube *Cube) {
	//fmt.Println(cube.location.x, cube.location.y, cube.location.z)
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				if x == 0 && y == 0 && z == 0 {
					continue
				}
				offset := int3{x, y, z}
				negOffset := int3{-x, -y, -z}
				neighborCoord := cube.location.Add(offset)

				if neighbor, ok := incomplete[neighborCoord]; ok {
					cube.setNeighbor(offset, neighbor)
					cube.presentNeighbors++

					//fmt.Println("Doing ", neighborCoord.x, neighborCoord.y, neighborCoord.z)
					neighbor.setNeighbor(negOffset, cube)
					neighbor.presentNeighbors++

					if neighbor.presentNeighbors == 26 {
						delete(incomplete, neighborCoord)
					}
				}
			}
		}
	}
	if cube.presentNeighbors < 26 {
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

	for y, line := range lines {
		for x, c := range line {
			if c == '#' {
				cube := Cube{on: true, location: int3{x, y, 0}}
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

	sceneDim := sceneMax.Sub(sceneMin).Add(int3{1, 1, 1})
	grid := make([]rune, sceneDim.x*sceneDim.y*sceneDim.z)
	for z := 0; z < sceneDim.z; z++ {
		for y := 0; y < sceneDim.y; y++ {
			for x := 0; x < sceneDim.x; x++ {
				pos := z*(sceneDim.x*sceneDim.y) + y*sceneDim.x + x
				grid[pos] = '.'
			}
		}
	}

	for _, cube := range allCubes {
		if cube.on {
			off := cube.location.Sub(sceneMin)
			pos := off.z*(sceneDim.x*sceneDim.y) + off.y*sceneDim.x + off.x
			grid[pos] = '#'
		}
	}

	for z := 1; z < sceneDim.z-1; z++ {
		for y := 1; y < sceneDim.y-1; y++ {
			for x := 1; x < sceneDim.x-1; x++ {
				pos := z*(sceneDim.x*sceneDim.y) + y*sceneDim.x + x
				fmt.Printf("%c", grid[pos])
			}
			fmt.Println("")
		}
		fmt.Println("")
	}

	onSum := 0
	for _, c := range allCubes {
		if c.on {
			onSum++
		}
	}
	fmt.Println(onSum)
}
