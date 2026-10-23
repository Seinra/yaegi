package main

func VS16[T any](vals ...T) T { return vals[0] }

func main() {
	s := []int{7, 8}
	println(VS16(s...))
}

// Output:
// 7
