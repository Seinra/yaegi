package main

import (
	"fmt"
	"hash/crc32"
)

// Lock for the typeAssert binary-concrete fallback (interp/run.go): a raw
// binary value (e.g. a hash from a binary call) asserted to an interpreted
// interface matches when its method set implements the target. The negative
// arm pins matcher precision: an unsatisfied method set must still panic,
// so the fallback cannot turn every assertion into a match.

type t13sizer interface {
	Size() int
	Reset()
}

type t13nomatch interface {
	Nope()
}

func mustPanic() (panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	h := crc32.NewIEEE()
	_ = h.(t13nomatch)
	return false
}

func main() {
	h := crc32.NewIEEE()
	s, ok := h.(t13sizer)
	if !ok {
		println("MISMATCH: no match")
		return
	}
	s.Reset()
	if s.Size() != 4 {
		println("MISMATCH: size")
		return
	}
	if _, ok := h.(t13nomatch); ok {
		println("MISMATCH: over-match")
		return
	}
	if !mustPanic() {
		println("MISMATCH: no panic")
		return
	}
	fmt.Println("t13assert0 ok:", s.Size())
}

// Output:
// t13assert0 ok: 4
