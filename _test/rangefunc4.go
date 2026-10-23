package main

import (
	"fmt"
	"slices"
)

// Keys mimics maps.Keys: it returns an iterator sequence over map keys.
func Keys[K comparable, V any](m map[K]V) func(yield func(K) bool) {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// Vals mimics slices.Values: it returns an iterator sequence over elements.
func Vals[E any](s []E) func(yield func(E) bool) {
	return func(yield func(E) bool) {
		for _, e := range s {
			if !yield(e) {
				return
			}
		}
	}
}

func main() {
	m := map[string]int{"b": 2, "a": 1, "c": 3}
	keys := []string{}
	for k := range Keys(m) {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	fmt.Println("keys:", keys)
	sum := 0
	for v := range Vals([]int{1, 2, 3}) {
		for w := range Vals([]int{10, 20}) {
			sum += v + w
		}
	}
	fmt.Println("nested:", sum)
}

// Output:
// keys: [a b c]
// nested: 102
