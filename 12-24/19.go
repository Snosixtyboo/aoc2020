package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Rule struct {
	raw     string
	options [][]int
}

func createRule(id2Rule map[int]Rule, id int) {
	updatedRule := id2Rule[id]
	updatedRule.options = append(updatedRule.options, make([]int, 0))

	rd := strings.NewReader(updatedRule.raw)
	var subID int
	var text string
	currOption := 0

	for rd.Len() > 0 {
		if _, err := fmt.Fscanf(rd, "%d", &subID); err == nil {
			updatedRule.options[currOption] = append(updatedRule.options[currOption], subID)
		} else if _, err := fmt.Fscanf(rd, "|"); err == nil {
			updatedRule.options = append(updatedRule.options, make([]int, 0))
			currOption++
		} else if _, err := fmt.Fscanf(rd, "%s", &text); err == nil {
			updatedRule.options[currOption] = append(updatedRule.options[currOption], -1)
		} else {
			fmt.Println(updatedRule.raw)
		}
	}
	id2Rule[id] = updatedRule
}

func expandToFit(t byte, id2Rule map[int]Rule, id int, path []int, nextExpansions *[][]int, last bool) {

	rule := id2Rule[id]

	if len(rule.options) == 1 && rule.options[0][0] == -1 {
		if rule.raw == string(t) {
			if last {
				path = append(path, -1)
			}
			*nextExpansions = append(*nextExpansions, path)
		}
		return
	}

	for _, option := range rule.options {
		newPath := make([]int, len(option[1:]))
		copy(newPath, option[1:])
		newPath = append(newPath, path...)
		expandToFit(t, id2Rule, option[0], newPath, nextExpansions, last)
	}
}

func verify(id2Rule map[int]Rule, text []byte) bool {
	expansions := make([][]int, len(id2Rule[0].options))
	copy(expansions, id2Rule[0].options)
	var nextExpansions [][]int

	for i := 0; i < len(text); i++ {
		for _, exp := range expansions {
			if len(exp) == 0 {
				continue
			}
			expandToFit(text[i], id2Rule, exp[0], exp[1:], &nextExpansions, i == len(text)-1)
		}
		expansions, nextExpansions = nextExpansions, expansions[0:0]
	}

	for _, exp := range expansions {
		if len(exp) == 1 && exp[0] == -1 {
			return true
		}
	}
	return false
}

func main() {
	var fileName string
	var version2 bool
	flag.StringVar(&fileName, "file", "data/in18.txt", "Input file to use")
	flag.BoolVar(&version2, "v2", false, "Use task2 version")
	flag.Parse()

	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.ReplaceAll(string(content), "\"", ""), "\n")

	var l, id int
	var raw string

	id2Rule := make(map[int]Rule)
	for l = 0; l < len(lines) && lines[l] != ""; l++ {
		rd := bufio.NewReader(strings.NewReader(lines[l]))
		fmt.Fscanf(rd, "%d:", &id)
		raw, _ = rd.ReadString('\n')
		id2Rule[id] = Rule{raw: strings.TrimSpace(raw)}
		createRule(id2Rule, id)
	}

	allMatched := 0
	for _, line := range lines {
		if len(line) > 0 && verify(id2Rule, []byte(line)) {
			allMatched++
		}
	}
	fmt.Println(allMatched)
}
