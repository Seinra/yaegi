package main

import "github.com/traefik/yaegi/_test/t13box"

func main() {
	b := &t13box.Box[int]{head: 42}
	println(b.head)
}

// Output:
// 42
