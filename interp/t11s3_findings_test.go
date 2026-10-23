package interp

import (
	"reflect"
	"testing"
)

// Slice-3 pins (T11): coverage/behavior pins for the T4 reliability
// findings backlog. Each test pins current gc-identical behavior, except
// TestGenericFuncDotImportSuffixCollision which is a regression lock for
// the fixed gta suffix-skip misroute (was RED: undefined: Hello).

// (b) T4-004: Use must stay idempotent when called twice with the same
// exports (GenericFunc recompilation must not fail redeclaration).
func TestGenericFuncUseTwiceIdempotent(t *testing.T) {
	i := New(Options{})
	exp := Exports{
		"guthib.com/generic/generic": map[string]reflect.Value{
			"Hello": reflect.ValueOf(GenericFunc("func Hello[T comparable](v T) *T {\n\treturn &v\n}")),
		},
	}
	if err := i.Use(exp); err != nil {
		t.Fatal("first Use:", err)
	}
	if err := i.Use(exp); err != nil {
		t.Fatal("second Use:", err)
	}
	res, err := i.Eval("generic.Hello(3)")
	if err != nil {
		t.Fatal("Eval after double Use:", err)
	}
	if res.Elem().Interface() != 3 {
		t.Fatal("expected &(3), got", res)
	}
}

// (c) T4-005: dot-importing a binary package whose path merely shares a
// name suffix with the guest package (example.com/amain vs main) must
// still compile its generics. gc oracle: equivalent module program
// prints 3 with exit 0.
func TestGenericFuncDotImportSuffixCollision(t *testing.T) {
	i := New(Options{})
	err := i.Use(Exports{
		"example.com/amain/amain": map[string]reflect.Value{
			"Hello": reflect.ValueOf(GenericFunc("func Hello[T any](v T) { println(v) }")),
		},
	})
	if err != nil {
		t.Fatal("Use:", err)
	}
	_, err = i.Eval(`
import . "example.com/amain"
func main() { Hello(3) }
`)
	if err != nil {
		t.Fatal("dot-import collision Eval:", err)
	}
}

// (a) T4-003: setExec with a nil generator must leave a defined exec
// (silent nop fallback). The fallback is a defensive-only guard: every
// cfg path sets gen explicitly, so this pins the guard itself.
func TestSetExecNilGenFallback(t *testing.T) {
	i := New(Options{})
	n := &node{interp: i}
	setExec(n)
	if n.exec == nil {
		t.Fatal("exec still nil after setExec with nil gen")
	}
}
