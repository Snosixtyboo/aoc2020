package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
)

type Coord struct {
	x, y int
}

type Dir int

const (
	north Dir = 0
	east  Dir = 1
	south Dir = 2
	west  Dir = 3
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in12.txt", "Input file to use")
	flag.Parse()
	content, _ := ioutil.ReadFile(fileName)
	rd := bytes.NewReader(content)

	currDir := east
	currPos := Coord{0, 0}

	var command rune
	var commandValue int
	for {
		if _, err := fmt.Fscanf(rd, "%c%d", &command, &commandValue); err == io.EOF {
			break
		}
		switch command {
		case 'N':
			currPos.y += commandValue
		case 'E':
			currPos.x += commandValue
		case 'S':
			currPos.y -= commandValue
		case 'W':
			currPos.x -= commandValue
		case 'L':
			currDir = Dir((int(currDir) + 3*commandValue/90) % 4)
		case 'R':
			currDir = Dir((int(currDir) + commandValue/90) % 4)
		case 'F':
			currPos.x += commandValue * ((-int(currDir - 2)) % 2)
			currPos.y += commandValue * ((-int(currDir - 1)) % 2)
		}
	}
	absCoord := currPos
	if absCoord.x < 0 {
		absCoord.x = -absCoord.x
	}
	if absCoord.y < 0 {
		absCoord.y = -absCoord.y
	}
	fmt.Println(absCoord.x + absCoord.y)
}
