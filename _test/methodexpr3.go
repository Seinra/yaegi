package main

import "fmt"

type pinner struct{ v bool }

func (c *pinner) set(v bool) bool { c.v = v; return c.v }
func (c pinner) get() bool        { return c.v }

type pholder struct{ pinner }

func main() {
	o := pholder{pinner{v: false}}
	f := (*pholder).set
	fmt.Println(f(&o, true))
	fmt.Println(o.v)
	fmt.Println((*pholder).get(&o))
}

// Output:
// true
// true
// true
