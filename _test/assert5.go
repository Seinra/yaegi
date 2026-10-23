package main

import (
	"fmt"
	"os"
)

type appErr struct{ msg string }

func (e appErr) Error() string { return e.msg }

func main() {
	// A nil interface value never satisfies an assertion (gc: ok == false).
	var e0 error
	v0, ok0 := e0.(error)
	fmt.Println(v0 == nil, ok0)
	// A concrete value held in an interface asserts directly, without unwrap.
	var e1 error = appErr{"boom"}
	v1, ok1 := e1.(appErr)
	fmt.Println(ok1, v1.msg)
	// The pre-existing path: interface values holding binaries still unwrap.
	f, err := os.Open("definitely-not-here-12345")
	if f != nil {
		f.Close()
	}
	var e2 error = err
	v2, ok2 := e2.(*os.PathError)
	fmt.Println(ok2, v2 != nil)
	s, ok3 := e2.(fmt.Stringer)
	fmt.Println(ok3)
	_ = s
}

// Output:
// true false
// true boom
// true true
// false
