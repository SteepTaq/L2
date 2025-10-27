package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	input := []string{
		"пятак", "пятка", "тяпка",
		"листок", "слиток", "столик",
		"стол",
	}
	anagrams := findAnagrams(input)

	for key, value := range anagrams {
		fmt.Println(key, value)
	}
}

func findAnagrams(dict []string) map[string][]string {
	groups := make(map[string][]string)

	for _, word := range dict {
		lowerWord := strings.ToLower(word)

		runes := []rune(word)
		sort.Slice(runes, func(i, j int) bool {
			return runes[i] < runes[j]
		})
		signature := string(runes)

		groups[signature] = append(groups[signature], lowerWord)
	}

	result := make(map[string][]string)
	for _, group := range groups {
		if len(group) > 1 {
			key := group[0]

			sort.Strings(group)

			result[key] = group
		}
	}

	return result
}
