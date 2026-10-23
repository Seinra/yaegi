package main

import (
	"fmt"
	"maps"
	"slices"
)

func strCmp(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func main() {
	m := map[string]int{"b": 2, "a": 1, "c": 3}
	fmt.Println(slices.Sorted(maps.Keys(m)))
	fmt.Println(slices.SortedFunc(maps.Keys(m), strCmp))
	fmt.Println(slices.SortedStableFunc(maps.Keys(m), strCmp))
	seq := func(yield func(int) bool) {
		for _, v := range []int{3, 1, 2} {
			if !yield(v) {
				return
			}
		}
	}
	fmt.Println(slices.Sorted(seq))
	fmt.Println(slices.Sorted(maps.Keys(map[string]int{})))
	fmt.Println(slices.Collect(seq))
	fmt.Println(slices.AppendSeq([]int{0}, seq))
}

// Output:
// [a b c]
// [a b c]
// [a b c]
// [1 2 3]
// []
// [3 1 2]
// [0 3 1 2]
