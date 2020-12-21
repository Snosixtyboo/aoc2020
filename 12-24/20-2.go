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
	unmatched    [4]bool
	numUnmatched int
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
	piece.unmatched[1], piece.unmatched[3] = piece.unmatched[3], piece.unmatched[1]
	piece.border[0], piece.border[2] = reverseBorder(piece.border[0], piece.dim), reverseBorder(piece.border[2], piece.dim)
	pieces[p] = piece
}

func flipPieceV(p int) {
	piece := pieces[p]
	piece.content = flipTileV(piece.dim, piece.content)
	piece.border[0], piece.border[2] = piece.border[2], piece.border[0]
	piece.unmatched[0], piece.unmatched[2] = piece.unmatched[2], piece.unmatched[0]
	piece.border[1], piece.border[3] = reverseBorder(piece.border[1], piece.dim), reverseBorder(piece.border[3], piece.dim)
	pieces[p] = piece
}

func rotatePiece(p int) {
	piece := pieces[p]
	piece.content = rotateTile(piece.dim, piece.content)

	var oldB [4]uint
	var oldM [4]bool
	for i := 0; i < 4; i++ {
		oldB[i] = piece.border[i]
		oldM[i] = piece.unmatched[i]
	}
	for i := 0; i < 4; i++ {
		piece.border[i] = oldB[(i+1)%4]
		if i == 1 || i == 3 {
			piece.border[i] = reverseBorder(piece.border[i], piece.dim)
		}
		piece.unmatched[i] = oldM[(i+1)%4]
	}
	pieces[p] = piece
}

func fillGrid(id int, x, y int) {
	if grid[y][x] != 0 {
		return
	}

	grid[y][x] = id
	piece := pieces[id] // left and down
	for c := 1; c <= 2; c++ {
		code := piece.border[c]
		hash := border2Hash(code, piece.dim)
		results := border2Tile[hash]

		if len(results) > 1 {
			var other int
			for entry := range results {
				if entry != id {
					other = entry
					break
				}
			}
			for border2Hash(pieces[other].border[(c+2)%4], pieces[other].dim) != hash {
				rotatePiece(other)
			}
			if pieces[other].border[(c+2)%4] != code {
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

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in20.txt", "Input file to use")
	flag.Parse()

	content, _ := ioutil.ReadFile(fileName)
	tiles := bytes.Split(content, []byte("\n\n"))

	gridDim := int(math.Sqrt(float64(len(tiles))))
	grid = make([][]int, gridDim)
	for g := 0; g < gridDim; g++ {
		grid[g] = make([]int, gridDim)
	}

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
					if key == border2Hash(b, p.dim) {
						p.unmatched[i] = true
					}
				}
				p.numUnmatched++
				pieces[id] = p
			}
		}
	}

	refPiece := -1
	for i, p := range pieces {
		if p.numUnmatched == 2 {
			refPiece = i
			break
		}
	}

	for rotation := 0; rotation < 4; rotation++ {
		if !pieces[refPiece].unmatched[1] && !pieces[refPiece].unmatched[2] {
			break
		}
		rotatePiece(refPiece)
	}
	fillGrid(refPiece, 0, 0)

	pieceDim := pieces[refPiece].dim
	imgDim := gridDim * (pieceDim - 2)
	outImg := make([][]byte, imgDim)
	for y := 0; y < gridDim; y++ {
		for ty := 1; ty < pieces[refPiece].dim-1; ty++ {
			outImg[y*(pieceDim-2)+ty-1] = make([]byte, imgDim)
			for x := 0; x < gridDim; x++ {
				for tx := 1; tx < pieces[refPiece].dim-1; tx++ {
					outImg[y*(pieceDim-2)+ty-1][x*(pieceDim-2)+tx-1] = pieces[grid[y][x]].content[ty][tx]
				}
			}
		}
	}

	monsterText := []string{
		"                  # ",
		"#    ##    ##    ###",
		" #  #  #  #  #  #   ",
	}

	allMonsters := 0
	for f := 0; f < 2 && allMonsters == 0; f++ {
		for r := 0; r < 4 && allMonsters == 0; r++ {
			checkImg := outImg
			if f > 0 {
				checkImg = flipTileH(imgDim, checkImg)
			}
			for y := 0; y < imgDim-len(monsterText); y++ {
				for x := 0; x < imgDim-len(monsterText[0]); x++ {
					present := true
					for my := 0; my < len(monsterText); my++ {
						for mx := 0; mx < len(monsterText[0]); mx++ {
							if monsterText[my][mx] == '#' && checkImg[y+my][x+mx] != '#' {
								present = false
								break
							}
						}
					}
					if present {
						allMonsters++
					}
				}
			}
			outImg = rotateTile(imgDim, outImg)
		}
	}
	wavesInMonster := 0
	waveCount := 0

	for y := 0; y < imgDim; y++ {
		waveCount += strings.Count(string(outImg[y]), "#")
	}
	for y := 0; y < len(monsterText); y++ {
		wavesInMonster += strings.Count(monsterText[y], "#")
	}
	fmt.Println(waveCount - wavesInMonster*allMonsters)
}
