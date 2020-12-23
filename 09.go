package main

import (
	"container/list"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	Version2, _ := strconv.ParseBool(os.Args[1])
	LookLength, _ := strconv.Atoi(os.Args[2])

	currPos := 0
	LookData := make([]int, LookLength)
	LookMap := make(map[int]int)

	bytes, _ := ioutil.ReadAll(os.Stdin)
	lines := strings.Split(string(bytes), "\n")
	for i, line := range lines {
		newNumber, _ := strconv.Atoi(line)

		success := false
		if i >= LookLength {
			for _, first := range LookData {
				target := newNumber - first
				neededCount := 1
				if target == first {
					neededCount = 2
				}
				if LookMap[target] >= neededCount {
					success = true
					break
				}
			}

			if !success {
				if !Version2 {
					log.Fatalf("XMAS is corrupt - can't be produced: %d\n", newNumber)
				}
				maintainList := list.New()
				currSum := 0
				for _, line := range lines {
					addNumber, _ := strconv.Atoi(line)
					currSum += addNumber
					maintainList.PushBack(addNumber)
					for currSum > newNumber {
						currSum -= maintainList.Front().Value.(int)
						maintainList.Remove(maintainList.Front())
					}
					if currSum == newNumber {
						min, max := math.MaxInt64, 0
						for p := maintainList.Front(); p != nil; p = p.Next() {
							if p.Value.(int) < min {
								min = p.Value.(int)
							}
							if p.Value.(int) > max {
								max = p.Value.(int)
							}
							log.Println(p.Value.(int))
						}
						log.Println("Success! Weakness =", min, "+", max, "=", min+max)
						return
					}
				}
			}
			evictNumber := LookData[currPos]
			LookMap[evictNumber]--
			if LookMap[evictNumber] == 0 {
				delete(LookMap, evictNumber)
			}
		}
		LookData[currPos] = newNumber
		LookMap[newNumber] = LookMap[newNumber] + 1
		currPos = (currPos + 1) % LookLength
	}
}
