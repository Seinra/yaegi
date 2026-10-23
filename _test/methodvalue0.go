package main

import "fmt"

type T struct{ N int }

func (t *T) Set(n int) { t.N = n }

func mk() T { return T{N: 1} }

func main() {
	f := mk().Set
	f(42)
	fmt.Println("callresult-ok")
}

// Error:
// cannot call pointer method Set on T
