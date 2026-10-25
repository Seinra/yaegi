package main

import "fmt"

// Lock for the isRegularCall wrap branch (interp/run.go call()): a nested
// call returning a CONCRETE value, fed to an interpreted func taking an
// interface arg, must be boxed into a valueInterface (the parameter slot
// carries the wrapper struct type; a raw struct is rejected by reflect).
// The boxed values are retained across loop iterations and asserted per
// iteration to an interface WIDER than the static source interface, which
// requires the concrete node to survive the trip. The genValueInterface
// snapshot discipline (interp/value.go) is characterized here too:
// retention stays stable with it in place (its load-bearing shape is
// driver-scale only; this probe stays GREEN with the snapshot disabled).

type t13item interface{ Name() string }

type t13box struct {
	N string
	T int
}

func (b t13box) Name() string { return b.N }

func (b t13box) Tag() int { return b.T }

func mkBoxC(n string) t13box { return t13box{N: n, T: len(n)} }

func identBox(v t13item) t13item { return v }

func main() {
	names := []string{"alfajor", "turron", "mate"}
	var kept []t13item
	for _, n := range names {
		kept = append(kept, identBox(mkBoxC(n)))
	}
	if len(kept) != 3 {
		println("MISMATCH len")
		return
	}
	for i, n := range names {
		if kept[i].Name() != n {
			println("MISMATCH kept")
			return
		}
		w, ok := kept[i].(interface {
			Name() string
			Tag() int
		})
		if !ok || w.Tag() != len(n) {
			println("MISMATCH assert")
			return
		}
	}
	fmt.Println("t13retain0 ok:", kept[0].Name(), kept[1].Name(), kept[2].Name())
}

// Output:
// t13retain0 ok: alfajor turron mate
