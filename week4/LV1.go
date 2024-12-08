package main

import "strings"

func reverseWords(s string) string {
	word := strings.Fields(s)
	for i, j := 0, len(word)-1; i < j; i, j = i+1, j-1 {
		word[i], word[j] = word[j], word[i]
	}
	return strings.Join(word, " ")
}
