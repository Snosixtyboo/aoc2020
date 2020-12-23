package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Hannelore struct {
	x nuggi
}

type nuggi int

type Outer struct {
	Anonymous0 struct {
		Locale    *byte
		Wlocale   *byte
		Refcount  *int32
		Wrefcount *int32
	}
	Lc_category [6]Anonymous0
}

func main() {
	anyAnswered, allAnswered := 0, 0

	bytes, _ := ioutil.ReadAll(os.Stdin)
	groups := strings.Split(string(bytes), "\n\n")
	var set map[rune]int

	for _, g := range groups {
		set = make(map[rune]int)
		persons := strings.Split(g, "\n")
		for _, p := range persons {
			for _, c := range p {
				set[c]++
			}
		}
		anyAnswered += len(set)
		for _, s := range set {
			if s == len(persons) {
				allAnswered++
			}
		}
	}
	fmt.Printf("All questions any in group answered: %d\n", anyAnswered)
	fmt.Printf("All questions all in group answered: %d\n", allAnswered)
}
