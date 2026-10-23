package main

import "fmt"

type emptyA struct{}
type emptyB struct{}

func main() {
	a := emptyA{}
	b := emptyB{}
	fmt.Println(a == emptyA(b))
	fmt.Println(emptyB(a) == b)
	c := emptyA(b)
	fmt.Println(c == a)
}

// Output:
// true
// true
// true
