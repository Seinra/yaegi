package main

import "fmt"

// Lock for the range-assign aliasing (interp/ast.go aRangeAssign plus the
// interp/cfg.go capture/alias trio) and the forStmt7 empty-body skip: with
// `for k, v = range ...` the iteration variables are the pre-existing
// ones, so post-loop reads observe the last iteration value, as in gc; and
// `for init; cond; post {}` with an empty body still wires cond/post, so a
// true condition does not end the function with zero values.

func main() {
	k := -1
	v := -1
	sum := 0
	for k, v = range []int{10, 20, 30} {
		sum += v
	}
	fmt.Println(k, v, sum)
	j := 0
	for j = 0; j < 3; j++ {
	}
	fmt.Println(j)
}

// Output:
// 2 30 60
// 3
