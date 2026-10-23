package main

import "fmt"

type Namer interface{ TableName(s string) string }

type fixedNamer struct{}

func (fixedNamer) TableName(s string) string { return "t_" + s }

type holder struct {
	Tag string
	Namer
}

func main() {
	h := holder{Tag: "x", Namer: fixedNamer{}}
	fmt.Println(h.TableName("a"))
	var n Namer = holder{Tag: "y", Namer: fixedNamer{}}
	fmt.Println(n.TableName("b"))
}

// Output:
// t_a
// t_b
