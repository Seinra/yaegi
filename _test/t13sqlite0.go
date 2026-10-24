package main

import (
	"database/sql"
	"database/sql/driver"
	"unsafe"
)

// Lock for the modernc.org/sqlite CFG wall (T2): an alias-defined struct
// whose field refers to an alias declared LATER in the file. GTA must
// revisit the incomplete alias declaration so the forward reference
// resolves; otherwise the field type stays a nilT placeholder and the
// chained selector below dies in lookupBinMethod2 with a CFG post-order
// panic (sqlite_windows_amd64.go:26076, FpInit on Lookaside).
type outer1 = struct {
	F Later
}

type Outer = outer1

type Later1 = struct {
	X uintptr
}

type Later = Later1

func fetch(p uintptr) uintptr {
	return (*Outer)(unsafe.Pointer(p)).F.X
}

var zarr = [9]int8{1, 2, 3, 4, 5, 6, 7, 8, 9}

var zptr = (*[9]int8)(unsafe.Pointer(&zarr[0]))

// Use before the definition below (as zOrd used at sqlite 56043 precedes
// its package-level definition at 56971): the declared type must be the
// dereferenced element type, not a pointer to it.
func zget(idx int) int8 { return zord[idx] }

var zord = *zptr

// t13driver exercises passing the result of an interpreted call to a binary
// interface parameter (go-sqlite init pattern: sql.Register(name, newDriver())).
// The call result must go through the interface wrapper, as for other
// argument shapes, else reflect rejects the method-less interpreted value.
type t13driver struct{}

func (d *t13driver) Open(name string) (driver.Conn, error) { return nil, nil }

func newT13Driver() *t13driver { return &t13driver{} }

// t13gopen exercises an interpreted call result flowing into an interpreted
// interface parameter, then asserted to an anonymous interface (gorm Open
// pattern: dialector.(interface{ Apply(*Config) error })). The return site
// wraps the value with the concrete type node; the call site must preserve
// it rather than clobber it with the callee signature node (which has no
// usable type), else the assertion crashes on a nil type.
type t13config struct{ K int }

type t13dialector interface{ Name() string }

type t13sqliteDial struct{ dsn string }

func (dial *t13sqliteDial) Name() string { return dial.dsn }
func (dial *t13sqliteDial) Apply(cfg *t13config) error {
	cfg.K = 7
	return nil
}

func t13open(dsn string) t13dialector { return &t13sqliteDial{dsn: dsn} }

func t13gopen(dial t13dialector) int {
	if applied, ok := dial.(interface{ Apply(*t13config) error }); ok {
		cfg := &t13config{}
		_ = applied.Apply(cfg)
		return cfg.K
	}
	return -1
}

// t13assignIface exercises assigning a concrete composite to an interface
// struct field through a pointer (gorm config.NamingStrategy pattern). The
// composite frame slot carries the struct type, so it must store the raw
// value and let readers wrap it in a valueInterface as needed.
type t13namer interface{ Name() string }

type t13namingStrategy struct{ MaxLen int }

func (strategy t13namingStrategy) Name() string { return "ns" }

type t13config2 struct{ NamingStrategy t13namer }

// addrDeref exercises &*p (address of a dereference), which denotes p itself
// (sqlite Xsqlite3Utf8Read pattern: &*(*uintptr)(unsafe.Pointer(pz))). A
// dereference is always addressable, even when the pointer comes from a
// call result.
func addrDeref() uintptr {
	var cell uintptr = 0x1234
	ptr := &*(*uintptr)(unsafe.Pointer(&cell))
	return *ptr
}

// spin exercises consecutive labels (ccgo shape "zulu_time:\n__10:\n..."):
// the inner label nests in the AST inside the outer labeled statement, so
// label pre-registration must descend the chain (sqlite parseDateOrTime).
func spin() int {
	i := 0
zulu:
__10:
	if !(i >= 2) {
		i++
		if i > 10 {
			goto zulu
		}
		goto __10
	}
	return i
}

func main() {
	outer := Outer{}
	outer.F.X = 42
	got := fetch(uintptr(unsafe.Pointer(&outer)))
	if got != 42 {
		println("MISMATCH:", got)
		return
	}
	println("t13sqlite0 ok:", got)
	if n := spin(); n != 2 {
		println("MISMATCH:", n)
		return
	}
	println("t13sqlite0 labels ok:", spin())
	if got := zget(1); got != 2 {
		println("MISMATCH:", got)
		return
	}
	println("t13sqlite0 deref ok:", zget(1))
	if got := addrDeref(); got != 0x1234 {
		println("MISMATCH:", got)
		return
	}
	println("t13sqlite0 addrderef ok:", addrDeref())
	sql.Register("t13sqlite0", newT13Driver())
	println("t13sqlite0 register ok")
	if got := t13gopen(t13open(":memory:")); got != 7 {
		println("MISMATCH:", got)
		return
	}
	println("t13sqlite0 gopen ok:", t13gopen(t13open(":memory:")))
	cfg2 := &t13config2{}
	cfg2.NamingStrategy = t13namingStrategy{MaxLen: 64}
	if cfg2.NamingStrategy.Name() != "ns" {
		println("MISMATCH")
		return
	}
	println("t13sqlite0 assigniface ok")
}

// Output:
// t13sqlite0 ok: 42
// t13sqlite0 labels ok: 2
// t13sqlite0 deref ok: 2
// t13sqlite0 addrderef ok: 4660
// t13sqlite0 register ok
// t13sqlite0 gopen ok: 7
// t13sqlite0 assigniface ok
