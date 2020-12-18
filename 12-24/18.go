package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func openLayer(N *int, values []int, ops []byte) {
	*N = *N + 1
	values[*N] = 0
	ops[*N] = '+'
}

func closeLayer(N *int, values []int, ops []byte) {
	values[*N-1] = applyOp(ops[*N-1], values[*N-1], values[*N])
	*N = *N - 1
}

func applyOp(op byte, a int, b int) int {
	res := 0
	switch op {
	case '+':
		res = a + b
	case '*':
		res = a * b
	default:
		log.Fatalf("Unexpected operator: %d", op)
	}
	return res
}

const STACKSIZE = 100

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

	var values [STACKSIZE]int
	var ops [STACKSIZE]byte
	totalSum := 0

	lines := strings.Split(strings.ReplaceAll(string(content), " ", ""), "\n")
	for _, line := range lines {

		N := 0
		ops[N] = '+'
		values[N] = 0

		rd := strings.NewReader(line)
		for rd.Len() > 0 {
			char, _ := rd.ReadByte()
			switch char {
			case '(':
				openLayer(&N, values[:], ops[:])
			case ')':
				closeLayer(&N, values[:], ops[:])
				for version2 && ops[N] == '*' {
					closeLayer(&N, values[:], ops[:])
				}
			case '+':
				ops[N] = char
			case '*':
				ops[N] = char
				if version2 {
					openLayer(&N, values[:], ops[:])
				}
			default:
				rd.UnreadByte()
				var val int
				_, err := fmt.Fscanf(rd, "%d", &val)
				if err != nil {
					log.Fatalf("Unexpected parsing error!")
				}
				values[N] = applyOp(ops[N], values[N], val)
			}
		}
		if version2 {
			for N > 0 {
				closeLayer(&N, values[:], ops[:])
			}
		} else if N != 0 {
			log.Fatal("Unexpected end depth:", N)
		}
		totalSum += values[N]
	}
	fmt.Println("Computed sum of results:", totalSum)
}
