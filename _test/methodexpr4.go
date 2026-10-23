package main

import "fmt"

type inner struct{ v int }

func (p *inner) set(x int) { p.v = x }
func (p inner) get() int   { return p.v }

type mid struct{ inner }

type outer struct {
	mid
	x   int
	arr [2]inner
}

func main() {
	o := outer{mid: mid{inner: inner{v: 1}}, x: 5, arr: [2]inner{{v: 2}, {v: 3}}}
	f := (*outer).set
	f(&o, 10)
	fmt.Println("outer:", o.get())
	g := (*inner).get
	fmt.Println("arr:", g(&o.arr[0]), g(&o.arr[1]))
	h := (*inner).set
	h(&o.arr[0], 20)
	fmt.Println("arr0after:", o.arr[0].get())
	fmt.Println("mid:", o.mid.get(), "x:", o.x, "arr1:", o.arr[1].get())
}

// Output:
// outer: 10
// arr: 2 3
// arr0after: 20
// mid: 10 x: 5 arr1: 3
