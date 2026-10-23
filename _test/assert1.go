package main

import (
	"fmt"
	"reflect"
	"time"
)

type TestStruct struct{}

func (t TestStruct) String() string {
	return "hello world"
}

type DummyStringer interface {
	String() string
}

// Types for concrete-target assertion with promoted pointer-receiver
// methods (issue T11-G1): ptrEmbed promotes *counter methods to its value
// method set, so n.(ptrEmbed) is a dynamic check, as in gc.
type Namer interface {
	TableName(s string) string
	IndexName(t, c string) string
}

type counter struct{ n int }

func (c *counter) TableName(s string) string    { return "t" }
func (c *counter) IndexName(t, s string) string { return "i" }

type ptrEmbed struct {
	Table string
	*counter
}

func main() {
	aType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	var t interface{}
	t = time.Nanosecond
	s, ok := t.(fmt.Stringer)
	if !ok {
		fmt.Println("time.Nanosecond does not implement fmt.Stringer")
		return
	}
	fmt.Println(s.String())
	fmt.Println(t.(fmt.Stringer).String())
	bType := reflect.TypeOf(time.Nanosecond)
	fmt.Println(bType.Implements(aType))

	// not redundant with the above, because it goes through a slightly different code path.
	if _, ok := t.(fmt.Stringer); !ok {
		fmt.Println("time.Nanosecond does not implement fmt.Stringer")
		return
	} else {
		fmt.Println("time.Nanosecond implements fmt.Stringer")
	}

	t = 42
	foo, ok := t.(fmt.Stringer)
	if !ok {
		fmt.Println("42 does not implement fmt.Stringer")
	} else {
		fmt.Println("42 implements fmt.Stringer")
		return
	}
	_ = foo

	if _, ok := t.(fmt.Stringer); !ok {
		fmt.Println("42 does not implement fmt.Stringer")
	} else {
		fmt.Println("42 implements fmt.Stringer")
		return
	}

	// TODO(mpl): restore when fixed
	// var tt interface{}
	var tt DummyStringer
	tt = TestStruct{}
	ss, ok := tt.(fmt.Stringer)
	if !ok {
		fmt.Println("TestStuct does not implement fmt.Stringer")
		return
	}
	fmt.Println(ss.String())
	fmt.Println(tt.(fmt.Stringer).String())

	if _, ok := tt.(fmt.Stringer); !ok {
		fmt.Println("TestStuct does not implement fmt.Stringer")
		return
	} else {
		fmt.Println("TestStuct implements fmt.Stringer")
	}

	var n Namer = &ptrEmbed{Table: "p", counter: &counter{}}
	if en, ok := n.(ptrEmbed); ok {
		fmt.Println("assert-ok", en.Table)
	} else {
		fmt.Println("assert-fail")
	}
	var i interface{} = ptrEmbed{Table: "q", counter: &counter{}}
	if e, ok := i.(ptrEmbed); ok {
		fmt.Println("same-ok", e.Table)
	} else {
		fmt.Println("same-fail")
	}
	e2 := i.(ptrEmbed)
	fmt.Println("single-ok", e2.Table)
	if p, ok := n.(*ptrEmbed); ok {
		fmt.Println("ptr-ok", p.Table)
	} else {
		fmt.Println("ptr-fail")
	}
}

// Output:
// 1ns
// 1ns
// true
// time.Nanosecond implements fmt.Stringer
// 42 does not implement fmt.Stringer
// 42 does not implement fmt.Stringer
// hello world
// hello world
// TestStuct implements fmt.Stringer
// assert-fail
// same-ok q
// single-ok q
// ptr-ok p
