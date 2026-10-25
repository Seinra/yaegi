package main

import "fmt"

type ctxt struct{ v bool }

func (c *ctxt) copy() bool { return c.v }

func (c ctxt) val() bool { return c.v }

func (c *ctxt) add(n int) int {
	if c.v {
		return n + 1
	}
	return n
}

func main() {
	c := &ctxt{v: true}
	var lower func(*ctxt) bool
	lower = (*ctxt).copy
	fmt.Println(lower(c))
	var get func(ctxt) bool
	get = ctxt.val
	fmt.Println(get(*c))
	fmt.Println((*ctxt).copy(c))
	fmt.Println((*ctxt).add(c, 41))
}

// Output:
// true
// true
// true
// 42
