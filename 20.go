package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/bits"
	"strings"
)

const PUZZLESIZE = 12

var attachment [PUZZLESIZE][PUZZLESIZE][4]uint
var grid [PUZZLESIZE][PUZZLESIZE]map[int]bool

var pieces = make(map[int]Piece)
var correctedPieces = make(map[int]Piece)

type Piece struct {
	dim          int
	content      [][]byte
	border       [4]uint
	unmatched    [4]bool
	numUnmatched int
	orientation  int
	flip         int
}

func (p Piece) orientedBorder(b int) uint {
	side := (b + p.orientation) % 4
	if ((p.flip%2 == 1) && side%2 == 0) || ((p.flip/2 == 1) && side%2 == 1) {
		side = (side + 2) % 4
	}
	return p.border[side]
}

func (p Piece) String() string {
	result := ""
	result += "Content:\n"
	for y := 0; y < p.dim; y++ {
		result += string(p.content[y]) + "\n"
	}
	result += "\n"
	return result
}

type Config struct {
	tile     int
	flippedH bool
	flippedV bool
	rotation int
}

type Connection struct {
	configA Config
	A       int
	configB Config
	B       int
}

func reverseBorder(border uint, dim int) uint {
	return bits.Reverse(border) >> (64 - dim)
}

func border2Hash(border uint, dim int) uint {
	hash := border
	testHash := reverseBorder(hash, dim)
	if testHash < hash {
		hash = testHash
	}
	return hash
}

var connections []Connection

var border2Tile = make(map[uint]map[int]bool)

func initWFC() {
	for y := 0; y < PUZZLESIZE; y++ {
		for x := 0; x < PUZZLESIZE; x++ {
			grid[y][x] = make(map[int]bool)
			for k, p := range pieces {
				if (x == 0 || x == PUZZLESIZE-1) && (y == 0 || y == PUZZLESIZE-1) { // corner
					if p.numUnmatched == 2 {
						grid[y][x][k] = true
					}
				} else if x == 0 || x == PUZZLESIZE-1 || y == 0 || y == PUZZLESIZE-1 { // border
					if p.numUnmatched == 1 {
						grid[y][x][k] = true
					}
				} else if x > 0 && x < PUZZLESIZE-1 && y > 0 && y < PUZZLESIZE-1 {
					if p.numUnmatched == 0 {
						grid[y][x][k] = true
					}
				}
			}
		}
	}
}

type Action struct {
	id int
	x  int
	y  int
}

func collapse(x, y int, id int, depth int) {
	wave := []Action{Action{id, x, y}}
	var newWave []Action

	for its := 0; len(wave) > 0 && its < 1000; its++ {
		for _, act := range wave {
			id = act.id
			x = act.x
			y = act.y

			piece := pieces[id]

			// Rotate/flip to fit
			sides := [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
			fit := false
		Loop:
			for f := 0; f < 4; f++ {
				for r := 0; r < 4; r++ {
					passed := 0
					for s := 0; s < 4; s++ {
						nx := x + sides[s][0]
						ny := y + sides[s][1]
						if nx < 0 || nx >= PUZZLESIZE || ny < 0 || ny >= PUZZLESIZE {
							passed++
							continue
						}
						border := transformBorder(piece, s, r, f)
						if attachment[y][x][s] != 0 {
							if border == attachment[y][x][s] {
								passed++
							}
						} else {
							for cand := range border2Tile[border2Hash(border, piece.dim)] {
								if ok, _ := grid[ny][nx][cand]; ok {
									passed++
									break
								}
							}
						}
					}
					if passed == 4 {
						piece.orientation = r
						piece.flip = f
						fit = true
						break Loop
					}
				}
			}

			for s := 0; s < 4; s++ {
				nx := x + sides[s][0]
				ny := y + sides[s][1]
				if nx < 0 || nx >= PUZZLESIZE || ny < 0 || ny >= PUZZLESIZE {
					continue
				}
				attachment[ny][nx][(s+2)%4] = transformBorder(piece, s, piece.orientation, piece.flip)
			}

			if !fit {
				log.Fatal("No fitting config found!")
			}

			pieces[id] = piece

			// Remove from all cells
			for i := 0; i < PUZZLESIZE; i++ {
				for j := 0; j < PUZZLESIZE; j++ {
					delete(grid[i][j], id)
				}
			}

			grid[y][x] = make(map[int]bool)
			grid[y][x][id] = true

			// Remove implausible neighbors
			for s := 0; s < 4; s++ {
				nx := x + sides[s][0]
				ny := y + sides[s][1]
				if nx < 0 || nx >= PUZZLESIZE || ny < 0 || ny >= PUZZLESIZE {
					continue
				}

				if len(grid[ny][nx]) == 1 {
					continue
				}

				newSuperPositions := make(map[int]bool)
				for id := range grid[ny][nx] {
					neighbor := pieces[id]
					for _, border := range neighbor.border {
						if border2Hash(border, piece.dim) == border2Hash(piece.orientedBorder(s), piece.dim) {
							newSuperPositions[id] = true
							break
						}
					}
				}
				grid[ny][nx] = newSuperPositions
				if len(newSuperPositions) == 1 {
					for key := range newSuperPositions {
						newWave = append(newWave, Action{x: nx, y: ny, id: key})
					}
				}
			}
		}
		wave, newWave = newWave, wave[0:0]
	}
}

func transformBorder(piece Piece, side int, rot int, flip int) uint {
	base := side
	rotated := (side + rot) % 4
	reverse := false

	if base/2 != rotated/2 {
		reverse = true
	}

	if (rotated%2 == 0 && flip%2 == 1) || (rotated%2 == 1 && flip/2 == 1) {
		rotated = (rotated + 2) % 4
	}
	border := piece.border[rotated]

	if reverse {
		border = reverseBorder(border, piece.dim)
	}
	if rotated%2 == 1 && flip%2 == 1 {
		border = reverseBorder(border, piece.dim)
	}
	if rotated%2 == 0 && flip/2 == 1 {
		border = reverseBorder(border, piece.dim)
	}
	return border
}

func printGridSuperPositions() {
	for y := 0; y < PUZZLESIZE; y++ {
		for x := 0; x < PUZZLESIZE; x++ {
			fmt.Printf("%3d ", len(grid[y][x]))
		}
		fmt.Println("")
	}
}

func rotateTile(dim int, content [][]byte) [][]byte {
	newContent := make([][]byte, len(content))
	for y := 0; y < dim; y++ {
		newContent[y] = make([]byte, len(content[y]))
		for x := 0; x < dim; x++ {
			newContent[y][x] = content[x][dim-1-y]
		}
	}
	return newContent
}

func flipTileV(dim int, content [][]byte) [][]byte {
	newContent := make([][]byte, len(content))
	for y := 0; y < dim; y++ {
		newContent[y] = make([]byte, len(content[y]))
		for x := 0; x < dim; x++ {
			newContent[y][x] = content[dim-1-y][x]
		}
	}
	return newContent
}

func flipTileH(dim int, content [][]byte) [][]byte {
	newContent := make([][]byte, len(content))
	for y := 0; y < dim; y++ {
		newContent[y] = make([]byte, len(content[y]))
		for x := 0; x < dim; x++ {
			newContent[y][x] = content[y][dim-1-x]
		}
	}
	return newContent
}

func correctTiles() {
	for id, p := range pieces {
		if p.flip/2 == 1 {
			p.content = flipTileH(p.dim, p.content)
		}
		if p.flip%2 == 1 {
			p.content = flipTileV(p.dim, p.content)
		}
		for r := 0; r < p.orientation; r++ {
			p.content = rotateTile(p.dim, p.content)
		}
		correctedPieces[id] = p
	}
}

func printPicture() {
	for y := 0; y < PUZZLESIZE; y++ {
		for ty := 9; ty >= 0; ty-- {
			for x := 0; x < PUZZLESIZE; x++ {
				for tx := 0; tx < 10; tx++ {
					for p := range grid[y][x] {
						c := correctedPieces[p].content[ty][tx]
						if tx == 0 || tx == 9 || ty == 0 || ty == 9 {
							fmt.Print(string(c))
						} else {
							fmt.Print(" ")
						}
					}
				}
			}
			fmt.Println("")
		}
	}
}

func main() {
	var fileName string
	var version2 bool
	flag.StringVar(&fileName, "file", "data/in20.txt", "Input file to use")
	flag.BoolVar(&version2, "v2", false, "Use task2 version")
	flag.Parse()

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	tiles := bytes.Split(content, []byte("\n\n"))

	for _, tile := range tiles {
		lines := bytes.Split(tile, []byte("\n"))

		piece := Piece{dim: len(lines[1]), content: lines[1:]}
		rd := bytes.NewReader(lines[0])
		var id int
		fmt.Fscanf(rd, "Tile %d:", &id)

		for i := 0; i < 4; i++ {
			down := i % 2
			off := 0
			if i == 1 || i == 2 {
				off = 1
			}

			x, y := down*off*(piece.dim-1), (down^1)*off*(piece.dim-1)
			piece.border[i] = 0
			for r := 0; r < piece.dim; r++ {
				if piece.content[y][x] == '#' {
					piece.border[i] |= (1 << r)
				}
				y += down
				x += down ^ 1
			}
		}
		pieces[id] = piece
	}

	for id, p := range pieces {
		for _, border := range p.border {
			hash := border2Hash(border, p.dim)
			if border2Tile[hash] == nil {
				border2Tile[hash] = make(map[int]bool)
			}
			border2Tile[hash][id] = true
		}
	}
	for key, ids := range border2Tile {
		if len(ids) == 1 {
			for id := range ids {
				p := pieces[id]
				for i, b := range p.border {
					p.unmatched[i] = key == border2Hash(b, p.dim)
				}
				p.numUnmatched++
				pieces[id] = p
			}
		}
	}

	cornerID := math.MaxInt64
	cornerProduct := 1

	fmt.Printf("%d", cornerProduct)
	for id, p := range pieces {
		if p.numUnmatched == 2 {
			if id < cornerID {
				cornerID = id
			}

			cornerProduct *= id
			fmt.Printf(" * %d", id)
		}
	}
	fmt.Println("\nProduct of corner piece IDs:", cornerProduct)

	initWFC()
	collapse(0, 0, cornerID, 0)
	correctTiles()

	monsterDimX := 20
	monsterDimY := 3
	monsterText := `                  # #    ##    ##    ### #  #  #  #  #  #   `

	finalImage := make([][]byte, 10*PUZZLESIZE)
	for y := 0; y < PUZZLESIZE; y++ {
		for ty := 0; ty < 10; ty++ {
			finalImage[y*10+ty] = make([]byte, 10*PUZZLESIZE)
			for x := 0; x < PUZZLESIZE; x++ {
				for tx := 0; tx < 10; tx++ {
					for p := range grid[y][x] {
						finalImage[y*10+ty][x*10+tx] = correctedPieces[p].content[9-ty][tx]
					}
				}
			}
		}
	}

	finalImage = make([][]byte, 8*PUZZLESIZE)
	for y := 0; y < PUZZLESIZE; y++ {
		for ty := 1; ty < 9; ty++ {
			finalImage[y*8+ty-1] = make([]byte, 8*PUZZLESIZE)
			for x := 0; x < PUZZLESIZE; x++ {
				for tx := 1; tx < 9; tx++ {
					for p := range grid[y][x] {
						finalImage[y*8+ty-1][x*8+tx-1] = correctedPieces[p].content[9-ty][tx]
					}
				}
			}
		}
	}

	allMonsters := 0
	for f := 0; f < 4 && allMonsters == 0; f++ {
		for r := 0; r < 4 && allMonsters == 0; r++ {
			checkImage := finalImage
			if f%2 == 1 {
				checkImage = flipTileH(8*PUZZLESIZE, checkImage)
			}
			if f/2 == 1 {
				checkImage = flipTileV(8*PUZZLESIZE, checkImage)
			}
			for y := 0; y < len(checkImage)-monsterDimY; y++ {
				for x := 0; x < len(checkImage[0])-monsterDimX; x++ {
					present := true
				Loop:
					for ty := 0; ty < monsterDimY; ty++ {
						for tx := 0; tx < monsterDimX; tx++ {
							if monsterText[ty*monsterDimX+tx] == '#' && checkImage[y+ty][x+tx] != '#' {
								present = false
								break Loop
							}
						}
					}
					if present {
						allMonsters++
					}
				}
			}
			finalImage = rotateTile(8*PUZZLESIZE, finalImage)
		}
	}
	allWaves := 0
	for _, row := range finalImage {
		allWaves += strings.Count(string(row), "#")
	}
	monsterWaves := strings.Count(monsterText, "#")
	fmt.Println("Turbulence in water:", allWaves-allMonsters*monsterWaves)
}
