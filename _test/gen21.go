package main

import (
	"fmt"
	"math"
)

func first[T any](s []T) T { return s[0] }

func main() {
	fmt.Println(first[math.Pi]([]float64{1.5}))
}

// Error:
// math.Pi is not a type
