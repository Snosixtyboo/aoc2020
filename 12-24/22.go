package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in22.txt", "Input file to use")
	flag.Parse()

	content, _ := ioutil.ReadFile(fileName)
	decks := bytes.Split(content, []byte("\n\n"))
	deckStr := [2][][]byte{bytes.Split(decks[0], []byte("\n"))[1:], bytes.Split(decks[1], []byte("\n"))[1:]}
	deck := [2][]int{make([]int, len(deckStr[0])), make([]int, len(deckStr[1]))}

	winner, maxNum := 0, 0
	for i := 0; i < 2; i++ {
		for j := 0; j < len(deckStr[i]); j++ {
			deck[i][j], _ = strconv.Atoi(string(deckStr[i][j]))
			if deck[i][j] > maxNum {
				maxNum = deck[i][j]
				winner = i
			}
		}
	}
	loser := winner ^ 1

	for len(deck[loser]) > 0 {
		if deck[0][0] > deck[1][0] {
			deck[0] = append(deck[0], deck[0][0], deck[1][0])
		} else {
			deck[1] = append(deck[1], deck[1][0], deck[0][0])
		}
		deck[0], deck[1] = deck[0][1:], deck[1][1:]
	}

	winningProduct := 0
	for l := 0; l < len(deck[winner]); l++ {
		winningProduct += deck[winner][l] * (len(deck[winner]) - l)
	}
	fmt.Println(winningProduct)
}
