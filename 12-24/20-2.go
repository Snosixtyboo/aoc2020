package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/bits"
	"strings"
)

var grid [][]int
var pieces = make(map[int]Piece)
var border2Tile = make(map[uint]map[int]bool)

type Piece struct {
	dim          int
	content      [][]byte
	border       [4]uint
	numUnmatched int
}

func reverseBorder(border uint, dim int) uint {
	return bits.Reverse(border) >> (64 - dim)
}

func border2Hash(hash uint, dim int) uint {
	testHash := reverseBorder(hash, dim)
	if testHash < hash {
		hash = testHash
	}
	return hash
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

func flipPieceH(p int) {
	piece := pieces[p]
	piece.content = flipTileH(piece.dim, piece.content)
	piece.border[1], piece.border[3] = piece.border[3], piece.border[1]
	for i := 0; i < 4; i++ {
		piece.border[i] = reverseBorder(piece.border[i], piece.dim)
	}
	pieces[p] = piece
}

func flipPieceV(p int) {
	piece := pieces[p]
	piece.content = flipTileV(piece.dim, piece.content)
	piece.border[0], piece.border[2] = piece.border[2], piece.border[0]
	for i := 0; i < 4; i++ {
		piece.border[i] = reverseBorder(piece.border[i], piece.dim)
	}
	pieces[p] = piece
}

func rotatePiece(p int) {
	piece := pieces[p]
	piece.content = rotateTile(piece.dim, piece.content)
	piece.border[0], piece.border[1], piece.border[2], piece.border[3] = piece.border[1], piece.border[2], piece.border[3], piece.border[0]
	pieces[p] = piece
}

func prepareRefPiece() int {
	refPiece := -1
	for _, ids := range border2Tile {
		if len(ids) == 1 {
			for id := range ids {
				p := pieces[id]
				p.numUnmatched++
				if p.numUnmatched == 2 {
					refPiece = id
					break
				}
				pieces[id] = p
			}
		}
	}
	for len(border2Tile[border2Hash(pieces[refPiece].border[1], pieces[refPiece].dim)]) < 2 ||
		len(border2Tile[border2Hash(pieces[refPiece].border[2], pieces[refPiece].dim)]) < 2 {
		rotatePiece(refPiece)
	}
	return refPiece
}

func fillGrid(id int, x, y int) {
	if grid[y][x] != 0 {
		return
	}
	grid[y][x] = id
	piece := pieces[id]
	// right and down
	for c := 1; c <= 2; c++ {
		opposite := (c + 2) % 4
		code := piece.border[c]
		hash := border2Hash(code, piece.dim)
		// Find tile matching border code, transform to fit
		if len(border2Tile[hash]) > 1 {
			var other int
			for entry := range border2Tile[hash] {
				if entry != id {
					other = entry
				}
			}
			for border2Hash(pieces[other].border[opposite], pieces[other].dim) != hash {
				rotatePiece(other)
			}
			if pieces[other].border[opposite] == code {
				if c == 1 {
					flipPieceV(other)
				} else {
					flipPieceH(other)
				}
			}
			if c == 1 {
				fillGrid(other, x+1, y)
			} else {
				fillGrid(other, x, y+1)
			}
		}
	}
}

func countWaves(pieceDim int, gridDim int) int {
	imgDim := gridDim * (pieceDim - 2)
	outImg := make([][]byte, imgDim)
	for y := 0; y < gridDim; y++ {
		for ty := 1; ty < pieceDim-1; ty++ {
			outImg[y*(pieceDim-2)+ty-1] = make([]byte, imgDim)
			for x := 0; x < gridDim; x++ {
				for tx := 1; tx < pieceDim-1; tx++ {
					outImg[y*(pieceDim-2)+ty-1][x*(pieceDim-2)+tx-1] = pieces[grid[y][x]].content[ty][tx]
				}
			}
		}
	}
	// Monster to look for
	monsterText := []string{
		"                  # ",
		"#    ##    ##    ###",
		" #  #  #  #  #  #   "}
	// Iterate over image and count occurrences
	numMonsters := 0
	for cfg := 0; cfg < 8 && numMonsters == 0; cfg++ {
		checkImg := outImg
		if cfg > 4 {
			checkImg = flipTileH(imgDim, checkImg)
		}
		for y := 0; y < imgDim-len(monsterText); y++ {
		NewMonster:
			for x := 0; x < imgDim-len(monsterText[0]); x++ {
				for my := 0; my < len(monsterText); my++ {
					for mx := 0; mx < len(monsterText[0]); mx++ {
						if monsterText[my][mx] == '#' && checkImg[y+my][x+mx] != '#' {
							continue NewMonster
						}
					}
				}
				numMonsters++
			}
		}
		outImg = rotateTile(imgDim, outImg)
	}
	wavesInImage, wavesInMonster := 0, 0
	for y := 0; y < imgDim; y++ {
		wavesInImage += strings.Count(string(outImg[y]), "#")
	}
	for y := 0; y < len(monsterText); y++ {
		wavesInMonster += strings.Count(monsterText[y], "#")
	}
	return wavesInImage - wavesInMonster*numMonsters
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in20.txt", "Input file to use")
	flag.Parse()
	content, _ := ioutil.ReadFile(fileName)
	tiles := bytes.Split(content, []byte("\n\n"))
	// Create grid
	gridDim := int(math.Sqrt(float64(len(tiles))))
	grid = make([][]int, gridDim)
	for g := 0; g < gridDim; g++ {
		grid[g] = make([]int, gridDim)
	}
	// Scan tiles and compute border codes
	var id int
	for _, tile := range tiles {
		lines := bytes.Split(tile, []byte("\n"))
		fmt.Fscanf(bytes.NewReader(lines[0]), "Tile %d:", &id)
		piece := Piece{dim: len(lines[1]), content: lines[1:]}
		// Compute binary border codes (counter clock-wise)
		x, y, dx, dy := 0, 0, 1, 0
		for i := 0; i < 4; i++ {
			for r := 0; r < piece.dim; r++ {
				if piece.content[y][x] == '#' {
					piece.border[i] |= (1 << r)
				}
				x, y = x+dx, y+dy
			}
			x, y, dx, dy = x-dx, y-dy, -dy, dx
		}
		pieces[id] = piece
	}
	// Convert border code to transform invariant hash and find matches
	for id, p := range pieces {
		for _, border := range p.border {
			hash := border2Hash(border, p.dim)
			if border2Tile[hash] == nil {
				border2Tile[hash] = make(map[int]bool)
			}
			border2Tile[hash][id] = true
		}
	}
	// Find starting piece, fill grid and count waves and monsters
	fillGrid(prepareRefPiece(), 0, 0)
	fmt.Println(countWaves(pieces[id].dim, gridDim))
}
