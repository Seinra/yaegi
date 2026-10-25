package main

func VS15[T any](vals ...T) T { return vals[0] }

func main() { println(VS15(1, 2, 3)) }

// Output:
// 1
