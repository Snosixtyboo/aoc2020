package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strings"
)

const numBits = 10
const allSeats = 128 * 8

func array2Int(arr []int) (result int) {
	for i := range arr {
		result |= arr[i] << i
	}
	return
}

func countNumbersWithBitSet(num int, bit int) int {
	bitNum := 1 << bit
	m := num / bitNum
	numsWithBit := (m/2)*bitNum + ((m % 2) * (num - m*bitNum))
	return numsWithBit
}

func main() {
	maxID := 0
	minID := math.MaxInt32

	bytes, _ := ioutil.ReadAll(os.Stdin)
	lines := strings.Split(string(bytes), "\n")

	var bitsRemaining [numBits]int
	for bit := range bitsRemaining {
		bitsRemaining[i] = countNumbersWithBitSet(allSeats, bit)
	}

	for _, line := range lines {
		id := 0
		reg, _ := regexp.Compile("B|R")
		matches := reg.FindAllStringIndex(strings.TrimSpace(line), -1)
		for _, pos := range matches {
			bit := (numBits - 1 - pos[0])
			id |= 1 << bit
			bitsRemaining[bit]--
		}
		if id > maxID {
			maxID = id
		}
		if id < minID {
			minID = id
		}
	}
	// Fill in the bits for seats in rows that are not available on plane
	for bit := 0; bit < numBits; bit++ {
		bitsRemaining[bit] -= countNumbersWithBitSet(minID, bit)
		bitsRemaining[bit] -= countNumbersWithBitSet(allSeats, bit) - countNumbersWithBitSet(maxID+1, bit)
	}
	// Find which one is actually mine by checking the only remaining bits
	myID := array2Int(bitsRemaining[:])
	fmt.Printf("The highest ID seen is %d\n", maxID)
	fmt.Printf("My ID is %d\n", myID)
}
