package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

const numBits = 10
const allSeats = 128 * 8

func main() {
	maxID := 0
	myID := 0

	bytes, _ := ioutil.ReadAll(os.Stdin)
	lines := strings.Split(string(bytes), "\n")

	var present [allSeats]bool

	for _, line := range lines {
		id := 0
		reg, _ := regexp.Compile("B|R")
		matches := reg.FindAllStringIndex(strings.TrimSpace(line), -1)
		for _, pos := range matches {
			bit := (numBits - 1 - pos[0])
			id |= 1 << bit
		}
		present[id] = true
		if id > maxID {
			maxID = id
		}
	}
	for id, taken := range present {
		if id > 0 && id < allSeats-1 {
			if !taken && present[id-1] && present[id+1] {
				myID = id
				break
			}
		}
	}
	fmt.Printf("The highest ID seen is %d\n", maxID)
	fmt.Printf("My ID is %d\n", myID)
}
