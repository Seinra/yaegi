package main

type props struct {
	pos  uint8
	size uint8
}

type input struct {
	str string
}

type form struct {
	name int
	info func(input, int) props
	next func(*iter) []byte
}

type buf struct {
	items [4]props
	f     form
	src   input
	flush func(*buf) bool
}

type iter struct {
	rb   buf
	data [8]byte
	next func(*iter) []byte
}

var anchor form

type writer struct {
	rb buf
}

var _ writer

func main() {
	anchor.name = 7
	println("rec-ok", anchor.name)
}

// Output:
// rec-ok 7
