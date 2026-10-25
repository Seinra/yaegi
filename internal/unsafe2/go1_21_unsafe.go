//go:build go1.21
// +build go1.21

// Package unsafe2 provides helpers to generate recursive struct types.
package unsafe2

import (
	"reflect"
	"unsafe"
)

type dummy struct{}

// DummyType represents a stand-in for a recursive type.
var DummyType = reflect.TypeOf(dummy{})

// The following type sizes must match their original definition in Go src/internal/abi/type.go.
type abiType struct {
	Size_      uintptr
	PtrBytes   uintptr
	Hash       uint32
	TFlag      uint8
	Align_     uint8
	FieldAlign_ uint8
	Kind_      uint8
	Equal      unsafe.Pointer
	GCData     *byte
	Str        int32
	PtrToThis  int32
}

type abiName struct {
	Bytes *byte
}

type abiStructField struct {
	Name   abiName
	Typ    *abiType
	Offset uintptr
}

type abiStructType struct {
	abiType
	PkgPath abiName
	Fields  []abiStructField
}

type emptyInterface struct {
	typ *abiType
	_   unsafe.Pointer
}

// SetFieldType sets the type of the struct field at the given index, to the given type.
//
// The struct type must have been created at runtime. This is very unsafe.
func SetFieldType(s reflect.Type, idx int, t reflect.Type) {
	if s.Kind() != reflect.Struct || idx >= s.NumField() {
		return
	}

	rtyp := unpackType(s)
	styp := (*abiStructType)(unsafe.Pointer(rtyp))
	f := styp.Fields[idx]
	f.Typ = unpackType(t)
	styp.Fields[idx] = f
}

func unpackType(t reflect.Type) *abiType {
	v := reflect.New(t).Elem().Interface()
	eface := *(*emptyInterface)(unsafe.Pointer(&v))
	return eface.typ
}

const (
	tFlagUncommon = 1 << 0
	tFlagNamed    = 1 << 2
)

type abiUncommonType struct {
	PkgPath int32
	Mcount  uint16
	Xcount  uint16
	Moff    uint32
	_       uint32
}

type namedStructType struct {
	abiStructType
	uncommon abiUncommonType
}

//go:linkname addReflectOff reflect.addReflectOff
func addReflectOff(ptr unsafe.Pointer) int32

func writeVarint(buf []byte, n int) int {
	for i := 0; ; i++ {
		b := byte(n & 0x7f)
		n >>= 7
		if n == 0 {
			buf[i] = b
			return i + 1
		}
		buf[i] = b | 0x80
	}
}

func isExported(name string) bool {
	if name == "" {
		return false
	}
	c := name[0]
	return c >= 'A' && c <= 'Z'
}

func newABIName(n, tag string, exported, embedded bool) abiName {
	if len(n) >= 1<<29 {
		panic("abi.NewName: name too long")
	}
	if len(tag) >= 1<<29 {
		panic("abi.NewName: tag too long")
	}
	var nameLen [10]byte
	var tagLen [10]byte
	nameLenLen := writeVarint(nameLen[:], len(n))
	tagLenLen := writeVarint(tagLen[:], len(tag))
	var bits byte
	l := 1 + nameLenLen + len(n)
	if exported {
		bits |= 1 << 0
	}
	if len(tag) > 0 {
		l += tagLenLen + len(tag)
		bits |= 1 << 1
	}
	if embedded {
		bits |= 1 << 3
	}
	b := make([]byte, l)
	b[0] = bits
	copy(b[1:], nameLen[:nameLenLen])
	copy(b[1+nameLenLen:], n)
	if len(tag) > 0 {
		tb := b[1+nameLenLen+len(n):]
		copy(tb, tagLen[:tagLenLen])
		copy(tb[tagLenLen:], tag)
	}
	return abiName{Bytes: &b[0]}
}

// NamedStructOf returns a named struct type with the given package path and name,
// based on the anonymous struct type anon (typically created by reflect.StructOf).
// It clones the type descriptor so the reflect.StructOf cache is not mutated,
// then sets the runtime name, package path and named/uncommon flags.
func NamedStructOf(anon reflect.Type, pkgPath, name string) reflect.Type {
	if anon.Kind() != reflect.Struct || pkgPath == "" || name == "" {
		return anon
	}
	orig := (*abiStructType)(unsafe.Pointer(unpackType(anon)))
	n := anon.NumField()
	newFields := make([]abiStructField, n)
	copy(newFields, orig.Fields)
	fullName := pkgPath + "." + name
	strName := newABIName(fullName, "", isExported(name), false)
	pkgName := newABIName(pkgPath, "", false, false)
	strOff := addReflectOff(unsafe.Pointer(strName.Bytes))
	pkgOff := addReflectOff(unsafe.Pointer(pkgName.Bytes))
	newDesc := &namedStructType{}
	newDesc.abiType = orig.abiType
	newDesc.PkgPath = orig.PkgPath
	newDesc.Fields = newFields
	newDesc.uncommon.Mcount = 0
	newDesc.uncommon.Xcount = 0
	newDesc.uncommon.Moff = uint32(unsafe.Sizeof(abiUncommonType{}))
	newDesc.uncommon.PkgPath = pkgOff
	newDesc.Str = strOff
	newDesc.TFlag |= tFlagUncommon | tFlagNamed
	type eface struct {
		typ  unsafe.Pointer
		data unsafe.Pointer
	}
	oldE := *(*eface)(unsafe.Pointer(&anon))
	newE := eface{typ: oldE.typ, data: unsafe.Pointer(newDesc)}
	return *(*reflect.Type)(unsafe.Pointer(&newE))
}
