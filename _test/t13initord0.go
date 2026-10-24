package main

import "fmt"

// Init-order: X depends on M transitively through reg body.
// gc initializes M before X (dependency order, not declaration order).
var X = reg("k")

func reg(k string) string {
	M[k] = 1
	return k
}

func main() { fmt.Println(X, M["k"]) }

var M = map[string]int{}

// Output:
// k 1
