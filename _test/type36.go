package main

import "fmt"

func first[T any](s []T) T { return s[0] }

func wrap[T any](v T) T { return first[T]([]T{v}) }

type box struct{ V int }

func main() {
	// Explicit type arguments which denote types keep working (no over-fire).
	fmt.Println(first[int]([]int{7, 8}))
	fmt.Println(wrap(42))
	fmt.Println(wrap(box{3}).V)
	// A local value shadowing nothing still infers normally.
	first := first[int]
	fmt.Println(first([]int{9}))
}

// Output:
// 7
// 42
// 3
// 9
