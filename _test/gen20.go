package main

import "fmt"

func first[T any](s []T) T { return s[0] }

func main() {
	fmt.Println(first[42]([]int{1, 2, 3}))
}

// Error:
// 42 is not a type
