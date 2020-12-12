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

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in12.txt", "Input file to use")
	flag.Parse()
	content, _ := ioutil.ReadFile(fileName)
	rd := bytes.NewReader(content)

	wayPos := Coord{10, 1}
	shipPos := Coord{0, 0}

	var command rune
	var commandValue int
	for {
		if _, err := fmt.Fscanf(rd, "%c%d\n", &command, &commandValue); err == io.EOF {
			break
		}
		switch command {
		case 'N':
			wayPos.y += commandValue
		case 'E':
			wayPos.x += commandValue
		case 'S':
			wayPos.y -= commandValue
		case 'W':
			wayPos.x -= commandValue
		case 'L':
			oldPos, quarterTurns, halfTurns := wayPos, commandValue/90, commandValue/180
			wayPos.x = ((quarterTurns+1)%2)*oldPos.x*(-2*(halfTurns%2)+1) + (quarterTurns%2)*oldPos.y*(2*(halfTurns%2)-1)
			wayPos.y = ((quarterTurns+1)%2)*oldPos.y*(-2*(halfTurns%2)+1) + (quarterTurns%2)*oldPos.x*(-2*(halfTurns%2)+1)
		case 'R':
			oldPos, quarterTurns, halfTurns := wayPos, commandValue/90, commandValue/180
			wayPos.x = ((quarterTurns+1)%2)*oldPos.x*(-2*(halfTurns%2)+1) + (quarterTurns%2)*oldPos.y*(-2*(halfTurns%2)+1)
			wayPos.y = ((quarterTurns+1)%2)*oldPos.y*(-2*(halfTurns%2)+1) + (quarterTurns%2)*oldPos.x*(2*(halfTurns%2)-1)
		case 'F':
			shipPos.x += commandValue * wayPos.x
			shipPos.y += commandValue * wayPos.y
		}
	}
	absCoord := shipPos
	if absCoord.x < 0 {
		absCoord.x = -absCoord.x
	}
	if absCoord.y < 0 {
		absCoord.y = -absCoord.y
	}
	fmt.Println(absCoord.x + absCoord.y)
}
