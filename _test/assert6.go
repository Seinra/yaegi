package main

import "fmt"

// Two-level embedded-pointer assertion (issue T11-G1 follow-up): methods
// promoted through two hops (pointer-then-pointer, value-then-pointer)
// stay dynamic checks, as in gc, matching the single-hop pin in assert1.go.
type Namer interface {
	TableName(s string) string
}

type inner struct{ n int }

func (c *inner) TableName(s string) string { return "t" + s }

type mid struct{ *inner }

type outer2 struct {
	Tag string
	*mid
}

type midV struct{ *inner }

type outerV struct {
	Tag string
	midV
}

func main() {
	var n Namer = &outer2{Tag: "p", mid: &mid{inner: &inner{}}}
	if e, ok := n.(outer2); ok {
		fmt.Println("ptrhop-ok", e.Tag)
	} else {
		fmt.Println("ptrhop-fail")
	}
	var m Namer = &outerV{Tag: "q", midV: midV{inner: &inner{}}}
	if e, ok := m.(outerV); ok {
		fmt.Println("valhop-ok", e.Tag)
	} else {
		fmt.Println("valhop-fail")
	}
	var i interface{} = outerV{Tag: "r", midV: midV{inner: &inner{}}}
	if e, ok := i.(outerV); ok {
		fmt.Println("same-ok", e.Tag)
	} else {
		fmt.Println("same-fail")
	}
}

// Output:
// ptrhop-fail
// valhop-fail
// same-ok r
