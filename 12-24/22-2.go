package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

func recGame(deckASrc []int, deckBSrc []int, depth int) bool {
	decks := [2][]int{make([]int, len(deckASrc)), make([]int, len(deckBSrc))}
	hands := [2]map[string]int{make(map[string]int), make(map[string]int)}
	copy(decks[0], deckASrc)
	copy(decks[1], deckBSrc)

	maxNum := 0
	maxDeck := 0
	for round := 1; len(decks[0]) > 0 && len(decks[1]) > 0; round++ {
		for d := 0; d < 2; d++ {
			code := ""
			for i := 0; i < len(decks[d]); i++ {
				code += strconv.Itoa(decks[d][i]) + ","
				if decks[d][i] > maxNum {
					maxNum = decks[d][i]
					maxDeck = d
				}
			}
			hands[d][code]++
			if hands[d][code] == 2 {
				return true
			}
		}
		if round == 1 && maxDeck == 0 && depth != 0 { //Player 1 will eventually win. Stole this from solution in reddit
			return true
		}
		drawnA, drawnB := decks[0][0], decks[1][0]
		decks[0], decks[1] = decks[0][1:], decks[1][1:]

		win1 := true
		if len(decks[0]) >= drawnA && len(decks[1]) >= drawnB {
			win1 = recGame(decks[0][:drawnA], decks[1][:drawnB], depth+1)
		} else if drawnB > drawnA {
			win1 = false
		}

		if win1 {
			decks[0] = append(decks[0], drawnA, drawnB)
		} else {
			decks[1] = append(decks[1], drawnB, drawnA)
		}
	}
	winningDeck := decks[0]
	if len(decks[0]) == 0 { // Player 2 won
		winningDeck = decks[1]
	}

	if depth == 0 {
		winningProduct := 0
		for l := 0; l < len(winningDeck); l++ {
			winningProduct += winningDeck[l] * (len(winningDeck) - l)
		}
		fmt.Println(winningProduct)
	}
	return len(decks[1]) == 0
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in22.txt", "Input file to use")
	flag.Parse()

	content, _ := ioutil.ReadFile(fileName)
	decks := bytes.Split(content, []byte("\n\n"))
	deckStr := [2][][]byte{bytes.Split(decks[0], []byte("\n"))[1:], bytes.Split(decks[1], []byte("\n"))[1:]}
	deck := [2][]int{make([]int, len(deckStr[0])), make([]int, len(deckStr[1]))}

	for i := 0; i < 2; i++ {
		for j := 0; j < len(deckStr[i]); j++ {
			deck[i][j], _ = strconv.Atoi(string(deckStr[i][j]))
		}
	}

	before := time.Now()
	recGame(deck[0], deck[1], 0)
	fmt.Println(time.Since(before))
}
