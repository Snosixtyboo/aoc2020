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

var allergen2Ingredient = make(map[string]map[string]bool)
var ingredient2Allergen = make(map[string]string)
var ingredientCounts = make(map[string]int)

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
			ingredientCounts[ing]++
		}

		for _, allergen := range allergens {
			compMap := make(map[string]bool)
			for _, ing := range ingredients {
				compMap[ing] = true
			}
			if allergen2Ingredient[allergen] == nil {
				allergen2Ingredient[allergen] = compMap
			}
			refMap := allergen2Ingredient[allergen]
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
			for v := range allergen2Ingredient[s] {
				ingredient2Allergen[v] = s
				for m := range allergen2Ingredient {
					if m == s {
						continue
					}
					if _, ok := allergen2Ingredient[m][v]; ok {
						delete(allergen2Ingredient[m], v)
						if len(allergen2Ingredient[m]) == 1 {
							newSingles = append(newSingles, m)
						}
					}
				}
			}
		}
		singles, newSingles = newSingles, singles[0:0]
	}

	sumNonAllergens := 0
	for ing, count := range ingredientCounts {
		if _, yes := ingredient2Allergen[ing]; !yes {
			sumNonAllergens += count
		}
	}
	fmt.Println("Appearances of non-allergens:", sumNonAllergens)

	var dangerous []string
	for k := range ingredient2Allergen {
		dangerous = append(dangerous, k)
	}
	sort.Slice(dangerous, func(i, j int) bool { return ingredient2Allergen[dangerous[i]] < ingredient2Allergen[dangerous[j]] })
	fmt.Print("Dangerous ingredient list:", strings.Join(dangerous, ","))
}
