package main

import "fmt"

// Pointer-to-array indexing auto-dereferences ((*p)[i]), which is always
// addressable, so method values with pointer receivers bind on array
// elements reached through a pointer, as in gc (native R3-ptr-array-index).
// Direct array indexing follows its base; both shapes mutate in place.
type E struct{ v int }

func (e *E) set(x int) { e.v = x }
func (e E) get() int   { return e.v }

func main() {
	arr := [2]E{}
	p := &arr
	g := p[1].set
	g(6)
	fmt.Println(arr[1].get())
	f := arr[0].set
	f(5)
	fmt.Println(arr[0].get())
	h := (&arr)[0].set
	h(7)
	fmt.Println(arr[0].get())
}

// Output:
// 6
// 5
// 7
