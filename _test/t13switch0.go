package main

import (
	"fmt"
	"hash"
	"hash/crc32"
	"time"
)

// Lock for the type-switch gate (interp/run.go _case): a wrapped
// (interpreted) value held in an interface slot matches concrete and
// interface case types by concrete-type identity (never by the wrapper
// name); a raw (binary) value matches its binary interface case; and a
// value matching nothing takes the no-match arm.

type t13swIface interface{ Name() string }

type t13swA struct{ N string }

func (a *t13swA) Name() string { return a.N }

type t13swB struct{ N string }

func (b t13swB) Name() string { return b.N }

type t13nomatch interface{ Nope() }

func mkSwA() t13swIface { return &t13swA{N: "wrapped"} }

func main() {
	n := 0
	// Wrapped value in an interface slot: single concrete case.
	e := interface{}(mkSwA())
	switch v := e.(type) {
	case *t13swA:
		if v.N != "wrapped" {
			println("MISMATCH concrete")
			return
		}
		n += 1
	default:
		println("MISMATCH concrete-default")
		return
	}
	// Wrapped value in an interface slot: single interface case.
	switch v := e.(type) {
	case t13swIface:
		if v.Name() != "wrapped" {
			println("MISMATCH iface")
			return
		}
		n += 10
	default:
		println("MISMATCH iface-default")
		return
	}
	// Wrapped value in an interface slot: multi concrete cases (v keeps
	// the guard type, so narrow it with an assertion).
	switch v := e.(type) {
	case *t13swB, *t13swA:
		a, ok := v.(*t13swA)
		if !ok || a.N != "wrapped" {
			println("MISMATCH multi")
			return
		}
		n += 100
	default:
		println("MISMATCH multi-default")
		return
	}
	// Raw binary value, binary interface case.
	h := crc32.NewIEEE()
	switch v := h.(type) {
	case hash.Hash:
		if v.Size() != 4 {
			println("MISMATCH raw-iface-size")
			return
		}
		n += 1000
	default:
		println("MISMATCH raw-iface")
		return
	}
	// Raw value via empty-interface conversion, concrete case.
	switch v := interface{}(time.Second).(type) {
	case time.Duration:
		if v != time.Second {
			println("MISMATCH raw-concrete")
			return
		}
		n += 10000
	default:
		println("MISMATCH raw-default")
		return
	}
	// Wrapped value, no case matches: must take default, never match on
	// the wrapper name itself.
	matched := false
	switch e.(type) {
	case *t13swB:
		matched = true
	case t13nomatch:
		matched = true
	default:
		n += 100000
	}
	if matched {
		println("MISMATCH wrapper-name")
		return
	}
	fmt.Println("t13switch0 ok:", n)
}

// Output:
// t13switch0 ok: 111111
