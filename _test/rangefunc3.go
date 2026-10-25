package main

import "fmt"

func seq1(yield func(int) bool) {
	for i := 1; i <= 4; i++ {
		if !yield(i) {
			return
		}
	}
}

func findFirst() int {
	for v := range seq1 {
		if v > 2 {
			return v
		}
	}
	return -1
}

func main() {
	total := 0
	for v := range seq1 {
		if v == 2 {
			continue
		}
		total += v
	}
	fmt.Println("cont:", total)
	c := 0
	for range seq1 {
		c++
	}
	fmt.Println("count:", c)
	fmt.Println("find:", findFirst())
	n := 0
	for v := range seq1 {
		if v == 3 {
			goto done
		}
		n += v
	}
done:
	fmt.Println("goto:", n)
}

// Output:
// cont: 8
// count: 4
// find: 3
// goto: 3
