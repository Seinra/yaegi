package main

import "fmt"

type Namer interface {
	TableName(s string) string
	IndexName(t, c string) string
}

type embeddedNamer struct {
	Table string
	Namer
}

type fixedNamer struct{}

func (fixedNamer) TableName(s string) string    { return "t_" + s }
func (fixedNamer) IndexName(t, c string) string { return "i_" + t }

func check(n Namer) string {
	if en, ok := n.(embeddedNamer); ok {
		return "embedded:" + en.Table
	}
	return "plain:" + n.TableName("x")
}

func main() {
	fmt.Println(check(embeddedNamer{Table: "p", Namer: fixedNamer{}}))
	fmt.Println(check(fixedNamer{}))
}

// Output:
// embedded:p
// plain:t_x
