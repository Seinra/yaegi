package main

import "fmt"

func gen() func(yield func(int) bool) {
	return func(yield func(int) bool) {
		for i := 1; i <= 3; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

func main() {
	for v := range gen() {
		fmt.Println(v)
	}
	v := 0
	for v = range gen() {
		fmt.Println("assign:", v)
	}
}

// Output:
// 1
// 2
// 3
// assign: 1
// assign: 2
// assign: 3
