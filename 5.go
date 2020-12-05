package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

var allSeats = 128 * 8

func main() {
	maxID := 0
	minID := math.MaxInt32

	bytes, _ := ioutil.ReadAll(os.Stdin)
	lines := strings.Split(string(bytes), "\n")

	var bitsRemaining [10]int
	for i := range bitsRemaining {
		bitsRemaining[i] = allSeats / 2
	}

	for _, line := range lines {
		id := 0
		for i, code := range line {
			bit := (9 - i)
			if code == 'B' || code == 'R' {
				id |= 1 << bit
				bitsRemaining[bit]--
			}
		}
		if id > maxID {
			maxID = id
		}
		if id < minID {
			minID = id
		}
	}
	// Fill in the bits for seats in rows that are not available on plane
	for i := 0; i < allSeats; i++ {
		if i == minID {
			i = maxID + 1
		}
		for bit := 0; bit < 10; bit++ {
			if i&(1<<bit) != 0 {
				bitsRemaining[bit]--
			}
		}
	}
	// Find which one is actually mine by checking the only remaining bits
	myID := 0
	for i := 0; i < 10; i++ {
		if bitsRemaining[i] != 0 {
			myID |= 1 << i
		}
	}
	fmt.Printf("The highest ID seen is %d\n", maxID)
	fmt.Printf("My ID is %d\n", myID)
}
