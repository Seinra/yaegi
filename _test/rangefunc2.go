package main

import "fmt"

func seqInf(yield func(int) bool) {
	i := 0
	for {
		i++
		if !yield(i) {
			return
		}
	}
}

type mySeq func(yield func(int) bool)

func genNamed() mySeq {
	return func(yield func(int) bool) {
		yield(42)
		yield(43)
	}
}

func main() {
	n := 0
	for v := range seqInf {
		n += v
		if v >= 4 {
			break
		}
	}
	fmt.Println("breaksum:", n)
	for v := range genNamed() {
		fmt.Println("named:", v)
	}
}

// Output:
// breaksum: 10
// named: 42
// named: 43
