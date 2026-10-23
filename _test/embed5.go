package main

import "fmt"

// Multi-embedded interface dispatch (issue T11-G2 follow-up): a holder
// embedding two interfaces dispatches each method through the interface
// which carries it, in field order. The first match no longer shadows
// the second (getMethodByName tries every exported valueInterface
// candidate until one matches); single-interface dispatch in embed4.go
// is unaffected.
type A interface{ Foo() string }
type B interface{ Bar() string }

type aImpl struct{}

func (aImpl) Foo() string { return "foo" }

type bImpl struct{}

func (bImpl) Bar() string { return "bar" }

type H struct {
	Tag string
	A
	B
}

func main() {
	h := H{Tag: "x", A: aImpl{}, B: bImpl{}}
	fmt.Println(h.Foo(), h.Bar(), h.Tag)
	var b B = h
	fmt.Println(b.Bar())
	var a A = h
	fmt.Println(a.Foo())
}

// Output:
// foo bar x
// bar
// foo
