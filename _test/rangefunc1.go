package main

import "fmt"

func seq2(yield func(int, string) bool) {
	vals := []string{"a", "b", "c"}
	for i, s := range vals {
		if !yield(i+1, s) {
			return
		}
	}
}

func main() {
	for k, v := range seq2 {
		fmt.Println(k, v)
	}
	for v := range seq2 {
		fmt.Println("first:", v)
	}
	for _, v := range seq2 {
		fmt.Println("second:", v)
	}
}

// Output:
// 1 a
// 2 b
// 3 c
// first: 1
// first: 2
// first: 3
// second: a
// second: b
// second: c
