package main

import "fmt"

func main() {
	var a, b int32
	for a = 1; a <= 9; a++ {
		for b = 1; b <= a; b++ { //确保只打印一半
			fmt.Printf("%d * %d = %d\t", b, a, a*b) // \t是制表符使对齐
			if a == b {
				fmt.Print("\n")
			}
		}
	}
}
