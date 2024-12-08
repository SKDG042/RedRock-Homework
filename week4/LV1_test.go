package main

import "testing"

func Test(t *testing.T) {
	word := "the sky is blue"
	result := reverseWords(word)
	if result != "blue is sky the" {
		t.Errorf("测试不通过，期待得到 %q, 但实际得到 %q", "blue is sky the", result)
	}
}
