package main

import (
	"github.com/sahilm/fuzzy"
)

func FoundRelatedItems(target string, sources []string) []string {
	var result []string
	matches := fuzzy.Find(target, sources)

	for _, match := range matches {
		result = append(result, sources[match.Index])
	}

	return result
}
