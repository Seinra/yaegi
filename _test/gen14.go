package main

import "fmt"

func first[T any](s []T) T { return s[0] }

var x = 3

func main() {
	fmt.Println(first[x]([]int{1, 2, 3}))
}

// Error:
// x is not a type
