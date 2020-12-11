package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var version2 bool
	var fileName string
	flag.BoolVar(&version2, "v2", false, "Use solution for second task")
	flag.StringVar(&fileName, "file", "data/in10.txt", "Input file to use")
	flag.Parse()
	bytes, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(bytes), "\n")
	nums, possibleWays := make([]int64, len(lines)+1), make([]int64, len(lines)+2)
	for i, line := range lines {
		nums[i+1], _ = strconv.ParseInt(line, 0, 64)
	}
	sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
	nums = append(nums, nums[len(nums)-1]+3)
	possibleWays[0] = 1

	differences := [4]int{}
	for i := 0; i < len(nums)-1; i++ {
		differences[nums[i+1]-nums[i]]++
		for n := i + 1; n < len(nums) && nums[n]-nums[i] <= 3; n++ {
			possibleWays[n] += possibleWays[i]
		}
	}
	if !version2 {
		fmt.Printf("Solution: %d\n", differences[1]*differences[3])
	} else {
		fmt.Printf("Solution: %d\n", possibleWays[len(possibleWays)-1])
	}
}
