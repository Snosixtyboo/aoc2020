package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type valRange struct {
	from, to int
}

type field struct {
	range1 valRange
	range2 valRange
}

type rangeList []valRange

func (r rangeList) Len() int {
	return len(r)
}
func (r rangeList) Swap(a, b int) {
	r[a], r[b] = r[b], r[a]
}
func (r rangeList) Less(a, b int) bool {
	return r[a].from < r[b].from
}

func main() {
	var fileName string
	var version2 bool
	flag.StringVar(&fileName, "file", "data/in16.txt", "Input file to use")
	flag.BoolVar(&version2, "v2", false, "Use task2 version")
	flag.Parse()

	content, _ := ioutil.ReadFile(fileName)
	lines := strings.Split(string(content), "\n")

	ticketFields := make(map[string]field)
	var valRanges rangeList

	var l int
	for l = 0; l < len(lines); l++ {
		fields := regexp.MustCompile("((?:.|\\s)+): (\\d+)-(\\d+) or (\\d+)-(\\d+)").FindStringSubmatch(lines[l])
		if len(fields) == 0 {
			break
		}
		var r1, r2 valRange
		r1.from, _ = strconv.Atoi(fields[2])
		r1.to, _ = strconv.Atoi(fields[3])
		r2.from, _ = strconv.Atoi(fields[4])
		r2.to, _ = strconv.Atoi(fields[5])
		valRanges = append(valRanges, r1, r2)

		ticketFields[fields[1]] = field{r1, r2}
	}

	cRanges := make(rangeList, 0, len(valRanges))
	cRanges = append(cRanges, valRanges[0])

	sort.Sort(valRanges)
	for _, r := range valRanges {
		if cRanges[len(cRanges)-1].to >= r.from {
			cRanges[len(cRanges)-1].to = r.to
		} else {
			cRanges = append(cRanges, r)
		}
	}

	numFields := len(ticketFields)

	for lines[l] != "your ticket:" {
		l++
	}

	myValueStrings := strings.Split(lines[l+1], ",")

	for lines[l] != "nearby tickets:" {
		l++
	}

	invalidSum := 0
	var validTickets [][]int

	for l = l + 1; l < len(lines); l++ {
		nums := make([]int, numFields)
		numberStrs := strings.Split(lines[l], ",")
		for n := 0; n < numFields; n++ {
			nums[n], _ = strconv.Atoi(numberStrs[n])
		}

		ticketValid := true
		for _, num := range nums {
			numValid := false
			for _, r := range cRanges {
				if num >= r.from && num <= r.to {
					numValid = true
					break
				}
			}
			if !numValid {
				invalidSum += num
				ticketValid = false
				if version2 {
					break
				}
			}
		}

		if ticketValid {
			validTickets = append(validTickets, nums)
		}
	}

	if !version2 {
		fmt.Println(invalidSum)
	}

	pos2Field := make(map[int][]string)

	for name, f := range ticketFields {
		for p := 0; p < numFields; p++ {
			possible := true
			for _, v := range validTickets {
				if (v[p] < f.range1.from || v[p] > f.range1.to) && (v[p] < f.range2.from || v[p] > f.range2.to) {
					possible = false
					break
				}
			}
			if possible {
				pos2Field[p] = append(pos2Field[p], name)
			}
		}
	}

	occCount := make(map[string]int)

	for _, p2f := range pos2Field {
		for _, f := range p2f {
			occCount[f]++
		}
	}

	single := ""
	for f, c := range occCount {
		if c == 1 {
			single = f
		}
	}

	field2Pos := make(map[string]int)
	for single != "" {
	Loop:
		for p, p2f := range pos2Field {
			for _, f := range p2f {
				if f == single {
					field2Pos[single] = p
					for _, f := range p2f {
						occCount[f]--
					}
					break Loop
				}
			}
		}
		delete(pos2Field, field2Pos[single])
		single = ""
		for f, c := range occCount {
			if c == 1 {
				single = f
			}
		}
	}

	product := 1
	for n, p := range field2Pos {
		if strings.HasPrefix(n, "departure") {
			val, _ := strconv.Atoi(myValueStrings[p])
			product *= val
		}
	}

	if version2 {
		fmt.Println(product)
	}
}
