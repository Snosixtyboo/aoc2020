package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

var all2Ing = make(map[string]map[string]bool)
var ingCount = make(map[string]int)
var ing2All = make(map[string]string)

type WordPair struct {
	native  string
	foreign string
}

type Dictionary []WordPair

func (wp Dictionary) Len() int {
	return len(wp)
}
func (wp Dictionary) Less(i, j int) bool {
	return wp[i].native < wp[j].native
}
func (wp Dictionary) Swap(i, j int) {
	wp[i], wp[j] = wp[j], wp[i]
}

func main() {
	var fileName string
	flag.StringVar(&fileName, "file", "data/in21.txt", "Input file to use")
	flag.Parse()

	content, _ := ioutil.ReadFile(fileName)
	lines := bytes.Split(content, []byte("\n"))

	var singles []string
	for _, line := range lines {
		regString := "((?:[a-z]+ )+)\\(contains ((?:(?:[a-z]+)(?:, )?)+)\\)"
		regResult := regexp.MustCompile(regString).FindStringSubmatch(string(line))
		ingredients := strings.Split(strings.TrimSpace(regResult[1]), " ")
		allergens := strings.Split(strings.Replace(regResult[2], ",", "", -1), " ")

		for _, ing := range ingredients {
			ingCount[ing]++
		}

		for _, allergen := range allergens {
			compMap := make(map[string]bool)
			for _, ing := range ingredients {
				compMap[ing] = true
			}
			if all2Ing[allergen] == nil {
				all2Ing[allergen] = compMap
				continue
			}
			refMap := all2Ing[allergen]
			for r := range refMap {
				if _, ok := compMap[r]; !ok {
					delete(refMap, r)
					if len(refMap) == 1 {
						singles = append(singles, allergen)
					}
				}
			}
		}
	}

	var newSingles []string
	for len(singles) > 0 {
		for _, s := range singles {
			for v := range all2Ing[s] {
				ing2All[v] = s
				for m := range all2Ing {
					if m == s {
						continue
					}
					if _, ok := all2Ing[m][v]; ok {
						delete(all2Ing[m], v)
						if len(all2Ing[m]) == 1 {
							newSingles = append(newSingles, m)
						}
					}
				}
			}
		}
		singles, newSingles = newSingles, singles[0:0]
	}

	sumNonAllergens := 0
	for ing, count := range ingCount {
		if _, yes := ing2All[ing]; !yes {
			sumNonAllergens += count
		}
	}
	fmt.Println("Appearances of non-allergens:", sumNonAllergens)

	var dict Dictionary
	for k, v := range ing2All {
		dict = append(dict, WordPair{native: v, foreign: k})
	}
	sort.Sort(dict)

	fmt.Print("Dangerous ingredient list: ")
	for i, entry := range dict {
		fmt.Printf("%s", entry.foreign)
		if i < len(dict)-1 {
			fmt.Printf(",")
		}
	}
}
