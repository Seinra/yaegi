package main

import "fmt"

type DB struct{}

func (db *DB) Unscoped() string { return "method" }

type Statement struct {
	*DB
	Unscoped bool
}

func main() {
	s := &Statement{DB: &DB{}, Unscoped: true}
	if s.Unscoped == true {
		fmt.Println("field:true")
	} else {
		fmt.Println("method-won")
	}
}

// Output:
// field:true
