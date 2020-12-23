package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type BusConstraint struct {
	interval int64
	when     int64
}

func gcd(a, b, rem int64) int64 {
	for b != rem {
		t := b
		b = a % b
		a = t
	}
	return a
}

func lcm(a, b int64) int64 {
	denom := gcd(a, b, 0)
	return (a / denom) * b
}

func main() {
	var fileName string
	var version2 bool
	flag.StringVar(&fileName, "file", "data/in13.txt", "Input file to use")
	flag.BoolVar(&version2, "v2", false, "Use task2 version")
	flag.Parse()
	content, _ := ioutil.ReadFile(fileName)
	rd := bytes.NewReader(content)
	var target int64
	var line string
	fmt.Fscanln(rd, &target)
	fmt.Fscanln(rd, &line)

	busNames := strings.Split(line, ",")

	var definiteBusses []BusConstraint
	unknownBusses := 0

	for i, busName := range busNames {
		busNum, err := strconv.ParseInt(busName, 0, 64)
		if err != nil {
			unknownBusses++
		} else {
			definiteBusses = append(definiteBusses, BusConstraint{busNum, int64(i)})
		}
	}

	toVal := definiteBusses[0].interval
	toInt := toVal

	for i := 1; i < len(definiteBusses); i++ {
		targetDiff := definiteBusses[i].when
		compInt := definiteBusses[i].interval
		for compNum := (toVal / compInt) * compInt; compNum-toVal != targetDiff; {
			if compNum-toVal > targetDiff {
				toVal += toInt
				compNum = (toVal / compInt) * compInt
			} else {
				compNum += compInt
			}
		}
		toInt = lcm(toInt, compInt)
	}
	fmt.Println(toVal)
}
