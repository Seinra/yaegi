package main

import (
	"math/bits"
	"unsafe"
)

const WSize = bits.UintSize / 8

func f() int { println("SIDE-EFFECT"); return 1 }

// Operand of unsafe.Sizeof is never evaluated (gc parity): only "1"
// may print, never "SIDE-EFFECT".
func main() { println(unsafe.Sizeof(f()) / WSize) }

// Output:
// 1
