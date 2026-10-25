package unsafe2_test

import (
	"reflect"
	"testing"

	"github.com/traefik/yaegi/internal/unsafe2"
)

func TestSwapFieldType(t *testing.T) {
	f := []reflect.StructField{
		{
			Name: "A",
			Type: reflect.TypeOf(int(0)),
		},
		{
			Name: "B",
			Type: reflect.PtrTo(unsafe2.DummyType),
		},
		{
			Name: "C",
			Type: reflect.TypeOf(int64(0)),
		},
	}
	typ := reflect.StructOf(f)
	ntyp := reflect.PtrTo(typ)

	unsafe2.SetFieldType(typ, 1, ntyp)

	if typ.Field(1).Type != ntyp {
		t.Fatalf("unexpected field type: want %s; got %s", ntyp, typ.Field(1).Type)
	}
}

func TestNamedStructOf(t *testing.T) {
	fields := []reflect.StructField{
		{Name: "ID", Type: reflect.TypeOf(uint(0))},
		{Name: "Nombre", Type: reflect.TypeOf("")},
	}
	anon := reflect.StructOf(fields)
	if anon.Name() != "" {
		t.Fatalf("anonymous struct must have empty name, got %q", anon.Name())
	}
	named := unsafe2.NamedStructOf(anon, "main", "Producto")
	if named.Name() != "Producto" {
		t.Fatalf("named struct name: want %q, got %q", "Producto", named.Name())
	}
	if named.PkgPath() != "main" {
		t.Fatalf("named struct pkg path: want %q, got %q", "main", named.PkgPath())
	}
	if named.Kind() != reflect.Struct {
		t.Fatalf("named struct kind: want struct, got %s", named.Kind())
	}
	if named.String() != "main.Producto" {
		t.Fatalf("named struct string: want %q, got %q", "main.Producto", named.String())
	}
	if anon.Name() != "" {
		t.Fatalf("StructOf cache mutated: anonymous name became %q", anon.Name())
	}
	other := unsafe2.NamedStructOf(anon, "main", "Otro")
	if other == named {
		t.Fatalf("distinct named types with same underlying must not be identical")
	}
	if other.Name() != "Otro" {
		t.Fatalf("second named struct name: want %q, got %q", "Otro", other.Name())
	}
	// Naming must keep recursive fixup working: SetFieldType on a named clone.
	recFields := []reflect.StructField{
		{Name: "Next", Type: unsafe2.DummyType},
	}
	recAnon := reflect.StructOf(recFields)
	recNamed := unsafe2.NamedStructOf(recAnon, "main", "Node")
	ptr := reflect.PtrTo(recNamed)
	unsafe2.SetFieldType(recNamed, 0, ptr)
	if recNamed.Field(0).Type.Elem() != recNamed {
		t.Fatalf("recursive fixup on named type failed")
	}
}
