package main

import "fmt"

type Namer interface {
	TableName(s string) string
	IndexName(t, c string) string
}

type counter struct{ n int }

func (c *counter) TableName(s string) string    { return fmt.Sprintf("t_%s_%d", s, c.n) }
func (c *counter) IndexName(t, s string) string { return "i_" + t }

type ptrEmbed struct {
	Table string
	*counter
}

func check(n Namer) string {
	if en, ok := n.(*ptrEmbed); ok {
		return "embedded:" + en.Table
	}
	return "plain:" + n.TableName("x")
}

func main() {
	fmt.Println(check(&ptrEmbed{Table: "p", counter: &counter{n: 3}}))
	fmt.Println(check(&counter{n: 4}))
}

// Output:
// embedded:p
// plain:t_x_4
