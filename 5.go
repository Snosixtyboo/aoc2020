package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
)

var allSeats = 128 * 8

func array2Int(arr []int) (result int) {
	for i := range arr {
		result |= arr[i] << i
	}
	return
}

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
		bitString := strings.NewReplacer("F", "0", "B", "1", "L", "0", "R", "1").Replace(strings.TrimSpace(line))
		for i, code := range bitString {
			bit := (9 - i)
			if code == '1' {
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
	myID := array2Int(bitsRemaining[:])
	fmt.Printf("The highest ID seen is %d\n", maxID)
	fmt.Printf("My ID is %d\n", myID)
}
