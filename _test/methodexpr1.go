package main

import "fmt"

type acc struct{ n int }

func (a *acc) sum(ns ...int) int {
	t := a.n
	for _, n := range ns {
		t += n
	}
	return t
}

type ctxt struct{ v bool }

func (c *ctxt) copy() bool { return c.v }

func main() {
	a := &acc{n: 1}
	f := (*acc).sum
	fmt.Println(f(a, 2, 3))
	fmt.Println((*acc).sum(a, 4, 5))

	func() {
		defer func() {
			if recover() != nil {
				fmt.Println("recovered")
			}
		}()
		fmt.Println((*ctxt).copy(nil))
		fmt.Println("no-panic")
	}()
}

// Output:
// 6
// 10
// recovered
