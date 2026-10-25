package main

import "unsafe"

// Mirrors modernc.org/libc RawMem ([1<<50-1]byte): slicing a dereferenced
// array pointer must not materialize the array in a function frame slot,
// which OOM-crashes the host at function entry via reflect.New.
type rawMem [1<<50 - 1]byte

func fill(s uintptr, v string) {
	b := (*(*rawMem)(unsafe.Pointer(s)))[:len(v):len(v)]
	copy(b, v)
}

func main() {
	buf := make([]byte, 8)
	fill(uintptr(unsafe.Pointer(&buf[0])), "alfajor!")
	println(string(buf))
}

// Output:
// alfajor!
