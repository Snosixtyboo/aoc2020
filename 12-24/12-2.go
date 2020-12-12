package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
)

type Coord struct { // Helper coord struct
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
			oldPos, qtTOdd, hemi := wayPos, (commandValue/90)%2, -2*((commandValue/180)%2)+1
			wayPos.x = (1-qtTOdd)*oldPos.x*(hemi) + qtTOdd*oldPos.y*(-hemi)
			wayPos.y = (1-qtTOdd)*oldPos.y*(hemi) + qtTOdd*oldPos.x*(hemi)
		case 'R':
			oldPos, qtTOdd, hemi := wayPos, (commandValue/90)%2, -2*((commandValue/180)%2)+1
			wayPos.x = (1-qtTOdd)*oldPos.x*(hemi) + qtTOdd*oldPos.y*(hemi)
			wayPos.y = (1-qtTOdd)*oldPos.y*(hemi) + qtTOdd*oldPos.x*(-hemi)
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
