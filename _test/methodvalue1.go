package main

import "fmt"

type T struct{ N int }

func (t *T) Set(n int) { t.N = n }
func (t T) Get() int   { return t.N }

func mk() T   { return T{N: 1} }
func mkp() *T { return &T{N: 2} }

func main() {
	v := T{N: 1}
	f := v.Set
	f(41)
	fmt.Println(v.N)
	p := mkp()
	g := p.Set
	g(5)
	fmt.Println(p.N)
	s := []T{{N: 1}}
	h := s[0].Set
	h(7)
	fmt.Println(s[0].N)
	fmt.Println(mk().Get())
	q := mk().Get
	fmt.Println(q())
}

// Output:
// 41
// 5
// 7
// 1
// 1
