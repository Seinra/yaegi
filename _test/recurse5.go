package main

import "fmt"

// T8 follow-up (R1/R2/R3): recursive types reached through every dummy
// capable kind — pointer, map, slice, chan and array of the recursive
// type — plus invocation (not just declaration) of recursive func-typed
// fields. Termination is pinned by completion itself: a non-terminating
// dummy fixup hangs this file instead of printing. No-redundancy (one
// rebuild per type) is interior to refType and not observable from guest
// syntax, so it is documented, not separately pinned.
type P struct {
	name string
	next *P
}

type M struct {
	m map[string]*M
}

type S struct {
	s []*S
}

type C struct {
	ch chan *C
}

type A struct {
	arr [2]*A
}

type F struct {
	apply func(*F) int
}

func main() {
	p := &P{name: "p"}
	p.next = &P{name: "q"}
	fmt.Println(p.next.name)
	m := &M{m: map[string]*M{"k": {}}}
	fmt.Println(len(m.m))
	s := &S{s: []*S{{}, {}}}
	fmt.Println(len(s.s))
	c := &C{ch: make(chan *C, 1)}
	c.ch <- &C{}
	fmt.Println(len(c.ch))
	a := &A{arr: [2]*A{{}, {}}}
	fmt.Println(len(a.arr))
	f := &F{}
	f.apply = func(x *F) int { return 42 }
	fmt.Println(f.apply(f))
}

// Output:
// q
// 1
// 2
// 1
// 2
// 42
