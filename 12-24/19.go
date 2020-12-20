package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
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
			fmt.Println("Unknown token ", updatedRule.raw)
		}
	}
	id2Rule[id] = updatedRule
}

func verify(id2Rule map[int]Rule, text []byte, stack []int, depth int) bool {
	if len(stack) == 0 {
		return false
	}

	ruleID := stack[0]
	stack = stack[1:]
	rule := id2Rule[ruleID]

	if len(rule.options) == 1 && rule.options[0][0] == -1 {
		if strings.HasPrefix(string(text), rule.raw) {
			text = text[len(rule.raw):]
			if len(stack) == 0 && len(text) == 0 { // perfect parsed
				return true
			}
			return verify(id2Rule, text, stack, depth+1)
		}
		return false
	}

	for _, opt := range rule.options {
		newStack := make([]int, len(opt))
		copy(newStack, opt)
		newStack = append(newStack, stack...)
		if verify(id2Rule, text, newStack, depth+1) {
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

	before := time.Now()
	id2Rule := make(map[int]Rule)
	for l = 0; l < len(lines) && lines[l] != ""; l++ {
		rd := bufio.NewReader(strings.NewReader(lines[l]))
		fmt.Fscanf(rd, "%d:", &id)
		raw, _ = rd.ReadString('\n')
		id2Rule[id] = Rule{raw: strings.TrimSpace(raw)}
		createRule(id2Rule, id)
	}

	allMatched := 0
	for ; l < len(lines); l++ {
		line := lines[l]
		if len(line) > 0 {
			if verify(id2Rule, []byte(line), id2Rule[0].options[0], 0) {
				allMatched++
			}
		}
	}
	fmt.Println(time.Since(before))
	fmt.Println(allMatched)
}
