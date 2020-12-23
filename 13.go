package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

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
	minWaitTime := target
	var bestBus int64 = -1
	fmt.Fscanln(rd, &line)
	busNames := strings.Split(strings.ReplaceAll(line, ",x", ""), ",")
	for _, busName := range busNames {
		busNum, _ := strconv.ParseInt(busName, 0, 64)
		waitCycles := target / busNum
		if target%busNum != 0 {
			waitCycles++
		}
		waitTime := waitCycles*busNum - target

		if waitTime < minWaitTime {
			minWaitTime = waitTime
			bestBus = busNum
		}
	}
	fmt.Println(bestBus, minWaitTime)
	fmt.Println(minWaitTime * bestBus)
}
