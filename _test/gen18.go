package main

func CF18[T any](a, b T) T { return a }

func main() { println(CF18(1, "s")) }

// Error:
// cannot convert untyped string to untyped int
