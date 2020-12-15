package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	var fileName string
	var version2 bool
	flag.StringVar(&fileName, "file", "data/in15.txt", "Input file to use")
	flag.BoolVar(&version2, "v2", false, "Use task2 version")
	flag.Parse()

	content, _ := ioutil.ReadFile(fileName)
	numStrings := strings.Split(string(content), ",")

	num := 0
	seen := make(map[int]int)
	for i, numString := range numStrings {
		num, _ = strconv.Atoi(numString)
		seen[num] = i + 1
	}
	delete(seen, num)

	var target int
	if version2 {
		target = 30000000
	} else {
		target = 2020
	}

	for t := len(numStrings) + 1; t <= target; t++ {
		when, ok := seen[num]
		seen[num] = t - 1
		if !ok {
			num = 0
		} else {
			num = t - 1 - when
		}
	}
	fmt.Printf("%d", num)
}
