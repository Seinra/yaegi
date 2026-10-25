package main

import "fmt"

func seq1(yield func(int) bool) {
	for i := 1; i <= 4; i++ {
		if !yield(i) {
			return
		}
	}
}

func main() {
	total := 0
outer:
	for v := range seq1 {
		if v == 2 {
			continue outer
		}
		total += v
	}
	fmt.Println("cont:", total)
	n := 0
outer2:
	for v := range seq1 {
		n += v
		if v == 2 {
			break outer2
		}
	}
	fmt.Println("break:", n)
}

// Output:
// cont: 8
// break: 3
