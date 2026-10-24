package main

import (
	"database/sql/driver"
	"fmt"
)

// Lock for the assignFromCall binary-interface wrapper (interp/run.go):
// assigning an interpreted multi-return call result into a binary
// (non-empty, non-error) interface destination (the database/sql driver
// Rows pattern: r, err = newRows(...) into r driver.Rows). Without the
// wrapper, reflect rejects the method-less interpreted value when stored
// in the binary-typed destination slot, as for other interface-boundary
// shapes. The error result keeps the raw path (identity preserved).

type t13rows struct {
	items []string
	pos   int
}

func (r *t13rows) Columns() []string { return []string{"nombre"} }

func (r *t13rows) Close() error { return nil }

func (r *t13rows) Next(dest []driver.Value) error {
	if r.pos >= len(r.items) {
		return fmt.Errorf("no more rows")
	}
	dest[0] = r.items[r.pos]
	r.pos++
	return nil
}

func mkRows() (*t13rows, error) {
	return &t13rows{items: []string{"alfajor", "turron"}}, nil
}

func main() {
	var r driver.Rows
	var err error
	r, err = mkRows()
	if err != nil {
		println("MISMATCH: err")
		return
	}
	cols := r.Columns()
	if len(cols) != 1 || cols[0] != "nombre" {
		println("MISMATCH: cols")
		return
	}
	dest := make([]driver.Value, 1)
	var names []string
	for {
		if nerr := r.Next(dest); nerr != nil {
			break
		}
		names = append(names, dest[0].(string))
	}
	if len(names) != 2 || names[0] != "alfajor" || names[1] != "turron" {
		println("MISMATCH: names")
		return
	}
	fmt.Println("t13iface0 ok:", len(names), names[0], names[1])
}

// Output:
// t13iface0 ok: 2 alfajor turron
