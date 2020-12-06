package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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
	fmt.Printf("All questions any in group answered: %d", anyAnswered)
	fmt.Printf("All questions all in group answered: %d", allAnswered)
}
