package main

import "fmt"

type Tabler interface{ TableName(s string) string }
type Indexer interface{ IndexName(t, c string) string }

type fixedTabler struct{}

func (fixedTabler) TableName(s string) string { return "t_" + s }

type fixedIndexer struct{}

func (fixedIndexer) IndexName(t, c string) string { return "i_" + t }

type multi struct {
	Tag string
	Tabler
	Indexer
}

func checkTabler(n Tabler) string {
	if m, ok := n.(multi); ok {
		return "embedded:" + m.Tag
	}
	return "plain:" + n.TableName("x")
}

func main() {
	fmt.Println(checkTabler(multi{Tag: "m", Tabler: fixedTabler{}, Indexer: fixedIndexer{}}))
	fmt.Println(checkTabler(fixedTabler{}))
}

// Output:
// embedded:m
// plain:t_x
