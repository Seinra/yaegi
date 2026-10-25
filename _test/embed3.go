package main

import "fmt"

type Namer interface {
	TableName(s string) string
	IndexName(t, c string) string
}

type fixedNamer struct{}

func (fixedNamer) TableName(s string) string    { return "t_" + s }
func (fixedNamer) IndexName(t, c string) string { return "i_" + t }

type embeddedNamer struct {
	Table string
	Namer
}

func main() {
	defer func() {
		if recover() != nil {
			fmt.Println("recovered-nil")
		}
	}()
	var n Namer = embeddedNamer{Table: "p"}
	fmt.Println(n.TableName("x"))
}

// Output:
// recovered-nil
