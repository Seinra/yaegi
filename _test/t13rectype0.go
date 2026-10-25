package main

import "fmt"

// Lock for the recursive-type rebuild (interp/type.go refType): struct
// fields holding func(*T), []func(*T) and map[string]func(*T) resolve to
// the in-flight type instead of freezing a dummy, for self-recursive and
// mutually-recursive shapes. A frozen dummy fails when the field is set
// or called.

type t13node struct {
	match  func(*t13node) bool
	fns    []func(*t13node) bool
	byName map[string]func(*t13node) bool
	next   *t13node
	peer   *t13leaf
}

type t13leaf struct {
	match func(*t13leaf) bool
	back  *t13node
}

func main() {
	n := &t13node{
		byName: map[string]func(*t13node) bool{
			"alive": func(x *t13node) bool { return x != nil },
		},
	}
	n.match = func(x *t13node) bool { return x == n }
	n.fns = []func(*t13node) bool{func(x *t13node) bool { return x.next == nil }}
	o := &t13leaf{back: n}
	o.match = func(x *t13leaf) bool { return x.back == n }
	n.peer = o
	if !n.match(n) {
		println("MISMATCH func")
		return
	}
	if !n.fns[0](n) {
		println("MISMATCH slice")
		return
	}
	if !n.byName["alive"](n) {
		println("MISMATCH map")
		return
	}
	if !o.match(o) || n.peer != o {
		println("MISMATCH mutual")
		return
	}
	fmt.Println("t13rectype0 ok")
}

// Output:
// t13rectype0 ok
