package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

func printList(input list.List) {
	for l := input.Front(); l != nil; l = l.Next() {
		fmt.Print(l.Value.(int), " ")
	}
	fmt.Println()
}

func nextWrap(l list.List, e *list.Element) *list.Element {
	ret := e.Next()
	if ret == nil {
		ret = l.Front()
	}
	return ret
}

func main() {
	var fileName string
	var version2 bool
	flag.StringVar(&fileName, "file", "data/in23.txt", "Input file to use")
	flag.BoolVar(&version2, "v2", false, "Use task2 version")
	flag.Parse()

	numRounds := 100
	numElems := 9
	if version2 {
		numRounds = 10000000
		numElems = 1000000
	}

	content, _ := ioutil.ReadFile(fileName)
	num2Entry := make([]*list.Element, numElems+1)

	var cups list.List
	for i := 0; i < numElems; i++ {
		if i < len(content) {
			num, _ := strconv.Atoi(string(content[i]))
			num2Entry[num] = cups.PushBack(num)
		} else {
			num2Entry[i+1] = cups.PushBack(i + 1)
		}
	}

	currCup := cups.Front()
	for round := 0; round < numRounds; round++ {
		targetVal := ((currCup.Value.(int) + (numElems - 2)) % numElems) + 1
		begin3 := nextWrap(cups, currCup)
		end3 := begin3
		for step := 0; step < 3; step++ {
			if end3.Value.(int) == targetVal {
				targetVal = ((targetVal + (numElems - 2)) % numElems) + 1
				step, end3 = -1, begin3
				continue
			}
			end3 = nextWrap(cups, end3)
		}
		t := nextWrap(cups, num2Entry[targetVal])
		a := begin3
		for step := 0; step < 3; step++ {
			b := nextWrap(cups, a)
			cups.MoveBefore(a, t)
			a = b
		}
		currCup = nextWrap(cups, currCup)
	}
	if version2 {
		A := nextWrap(cups, num2Entry[1])
		B := nextWrap(cups, A)
		fmt.Println(A.Value.(int), B.Value.(int), A.Value.(int)*B.Value.(int))
	} else {
		for d := nextWrap(cups, num2Entry[1]); d != num2Entry[1]; d = nextWrap(cups, d) {
			fmt.Print(d.Value.(int))
		}
	}
}
