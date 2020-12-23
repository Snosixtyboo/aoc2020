package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var memory = make(map[uint64]uint64)

func allPermutations(addr uint64, xLocations [][]int) []uint64 {
	locations := make([]uint64, 0, 1<<len(xLocations))
	vLocations := make([]bool, len(xLocations))

	for i := 0; i < len(xLocations); i++ { // set all floating to 1 initially
		bit := 35 - xLocations[i][0]
		addr &^= (1 << bit)
	}

	setFloating := 0
	for setFloating != len(xLocations) {

		locations = append(locations, addr)

		for i := 0; i < len(vLocations); i++ {
			bit := 35 - xLocations[i][0]
			if !vLocations[i] {
				vLocations[i] = true
				addr |= (1 << bit)
				setFloating++
				break
			} else {
				vLocations[i] = false
				addr &^= (1 << bit)
				setFloating--
			}
		}
	}

	locations = append(locations, addr)
	return locations
}

func main() {
	var fileName string
	var version2 bool
	flag.StringVar(&fileName, "file", "data/in14.txt", "Input file to use")
	flag.BoolVar(&version2, "v2", false, "Use task2 version")
	flag.Parse()

	content, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(content), "\n")

	maskOne := uint64(0)
	maskZero := uint64(math.MaxUint64)
	var xLocations [][]int

	for _, line := range lines {
		rd := strings.NewReader(line)

		var instrStr, valStr string
		fmt.Fscanf(rd, "%s = %s", &instrStr, &valStr)
		instr := regexp.MustCompile("[a-zA-Z]+").FindString(instrStr)

		switch instr {
		case "mem":
			addr, _ := strconv.ParseUint(instrStr[4:len(instrStr)-1], 0, 64)
			val, _ := strconv.ParseUint(valStr, 0, 64)
			if version2 {
				addr = (addr | maskOne)
				for _, permAddr := range allPermutations(addr, xLocations) {
					memory[permAddr] = val
				}
			} else {
				val = (val | maskOne) & maskZero
				memory[addr] = val
			}
		case "mask":
			xLocations = regexp.MustCompile("X").FindAllStringIndex(valStr, -1)
			maskOne, _ = strconv.ParseUint(strings.ReplaceAll(valStr, "X", "0"), 2, 64)
			maskZero, _ = strconv.ParseUint(strings.ReplaceAll(valStr, "X", "1"), 2, 64)
		default:
			log.Fatalf("Unknown instruction!")
		}
	}
	sum := uint64(0)
	for _, val := range memory {
		sum += val
	}
	fmt.Println(sum)
}
