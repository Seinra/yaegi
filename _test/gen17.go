package main

func MK17[M ~map[K]V, K comparable, V any](m M) V {
	for _, v := range m {
		return v
	}
	var zero V
	return zero
}

func main() { println(MK17(map[string]int{"a": 42})) }

// Output:
// 42
