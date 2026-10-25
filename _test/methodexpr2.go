package main

import "fmt"

type inner struct{ v bool }

func (c inner) val() bool { return c.v }
func (c inner) add(n int) int {
	if c.v {
		return n + 1
	}
	return n
}

type outer struct{ inner }

func main() {
	o := outer{inner{v: true}}
	f := outer.val
	fmt.Println(f(o))
	fmt.Println(outer.val(o))
	fmt.Println(outer.add(o, 41))
}

// Output:
// true
// true
// 42
