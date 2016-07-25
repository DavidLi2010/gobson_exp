// Copyright 2015-2016 David Li
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package bson_test

import (
	"testing"

	"math"

	"reflect"

	"github.com/DavidLi2010/gobson_exp/bson"
)

type primary struct {
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uintptr uintptr
	Float32 float32
	Float64 float64
	String  string
}

type primaryPtr struct {
	Bool    *bool
	Int     *int
	Int8    *int8
	Int16   *int16
	Int32   *int32
	Int64   *int64
	Uint    *uint
	Uint8   *uint8
	Uint16  *uint16
	Uint32  *uint32
	Uint64  *uint64
	Uintptr *uintptr
	Float32 *float32
	Float64 *float64
	String  *string
}

var pridata = primary{
	Bool:    true,
	Int:     math.MinInt32,
	Int8:    math.MinInt8,
	Int16:   math.MinInt16,
	Int32:   math.MinInt32,
	Int64:   math.MinInt64,
	Uint:    math.MaxUint32,
	Uint8:   math.MaxUint8,
	Uint16:  math.MaxUint16,
	Uint32:  math.MaxUint32,
	Uint64:  math.MaxInt64,
	Uintptr: math.MaxInt32,
	Float32: math.MaxFloat32,
	Float64: math.MaxFloat64,
	String:  "string",
}

var mdata = map[string]interface{}{
	"Bool":    true,
	"Int":     int(math.MinInt32),
	"Int8":    int8(math.MinInt8),
	"Int16":   int16(math.MinInt16),
	"Int32":   int32(math.MinInt32),
	"Int64":   int64(math.MinInt64),
	"Uint":    uint(math.MaxUint32),
	"Uint8":   uint8(math.MaxUint8),
	"Uint16":  uint16(math.MaxUint16),
	"Uint32":  uint32(math.MaxUint32),
	"Uint64":  uint64(math.MaxInt64),
	"Uintptr": uintptr(math.MaxInt32),
	"Float32": float32(math.MaxFloat32),
	"Float64": float64(math.MaxFloat64),
	"String":  "string",
}

var mdata2 = bson.Map{
	"Bool":    true,
	"Int":     int(math.MinInt32),
	"Int8":    int8(math.MinInt8),
	"Int16":   int16(math.MinInt16),
	"Int32":   int32(math.MinInt32),
	"Int64":   int64(math.MinInt64),
	"Uint":    uint(math.MaxUint32),
	"Uint8":   uint8(math.MaxUint8),
	"Uint16":  uint16(math.MaxUint16),
	"Uint32":  uint32(math.MaxUint32),
	"Uint64":  uint64(math.MaxInt64),
	"Uintptr": uintptr(math.MaxInt32),
	"Float32": float32(math.MaxFloat32),
	"Float64": float64(math.MaxFloat64),
	"String":  "string",
}

var ddata = bson.Doc{
	{"Bool", true},
	{"Int", int(math.MinInt32)},
	{"Int8", int8(math.MinInt8)},
	{"Int16", int16(math.MinInt16)},
	{"Int32", int32(math.MinInt32)},
	{"Int64", int64(math.MinInt64)},
	{"Uint", uint(math.MaxUint32)},
	{"Uint8", uint8(math.MaxUint8)},
	{"Uint16", uint16(math.MaxUint16)},
	{"Uint32", uint32(math.MaxUint32)},
	{"Uint64", uint64(math.MaxInt64)},
	{"Uintptr", uintptr(math.MaxInt32)},
	{"Float32", float32(math.MaxFloat32)},
	{"Float64", float64(math.MaxFloat64)},
	{"String", "string"},
}

func valueEqual(v1, v2 interface{}) bool {
	var eq bool
	value1 := reflect.ValueOf(v1)
	value2 := reflect.ValueOf(v2)

	switch value1.Kind() {
	case reflect.String:
		fallthrough
	case reflect.Bool:
		eq = (v1 == v2)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch value2.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rv1 := value1.Int()
			rv2 := value2.Int()
			eq = (rv1 == rv2)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch value2.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rv1 := value1.Uint()
			rv2 := uint64(value2.Int())
			eq = (rv1 == rv2)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			rv1 := value1.Uint()
			rv2 := value2.Uint()
			eq = (rv1 == rv2)
		}
	case reflect.Float32, reflect.Float64:
		switch value2.Kind() {
		case reflect.Float32, reflect.Float64:
			rv1 := value1.Float()
			rv2 := value2.Float()
			eq = (rv1 == rv2)
		}
	case reflect.Ptr:
		switch value2.Kind() {
		case reflect.Ptr:
			eq = valueEqual(value1.Elem().Interface(), value2.Elem().Interface())
		}
	}

	return eq
}

func docValueEqual(d1, d2 bson.Doc) bool {
	if len(d1) != len(d2) {
		return false
	}

	for i := 0; i < len(d1); i++ {
		e := d1[i]
		e2 := d2[i]

		if e.Name != e2.Name {
			return false
		}

		if !valueEqual(e.Value, e2.Value) {
			return false
		}
	}

	return true
}

func TestPrimaryStruct(t *testing.T) {
	s := pridata

	b := bson.StructToBson(s)

	var s2 primary
	b.Struct(&s2)

	if s != s2 {
		t.Errorf("invalid primary type mapping")
	}
}

func TestPrimaryPointerStruct(t *testing.T) {
	Bool := true
	Int := int(math.MinInt32)
	Int8 := int8(math.MinInt8)
	Int16 := int16(math.MinInt16)
	Int32 := int32(math.MinInt32)
	Int64 := int64(math.MinInt64)
	Uint := uint(math.MaxUint32)
	Uint8 := uint8(math.MaxUint8)
	Uint16 := uint16(math.MaxUint16)
	Uint32 := uint32(math.MaxUint32)
	Uint64 := uint64(math.MaxInt64)
	Uintptr := uintptr(math.MaxInt32)
	Float32 := float32(math.MaxFloat32)
	Float64 := float64(math.MaxFloat64)
	String := "string"

	p := primaryPtr{
		Bool:    &Bool,
		Int:     &Int,
		Int8:    &Int8,
		Int16:   &Int16,
		Int32:   &Int32,
		Int64:   &Int64,
		Uint:    &Uint,
		Uint8:   &Uint8,
		Uint16:  &Uint16,
		Uint32:  &Uint32,
		Uint64:  &Uint64,
		Uintptr: &Uintptr,
		Float32: &Float32,
		Float64: &Float64,
		String:  &String,
	}

	var p2 primaryPtr
	bson.StructToBson(p).Struct(&p2)

	value := reflect.ValueOf(p)
	value2 := reflect.ValueOf(p2)
	typ := reflect.TypeOf(primaryPtr{})
	for i := 0; i < typ.NumField(); i++ {
		f := value.Field(i).Interface()
		f2 := value2.Field(i).Interface()

		if !valueEqual(f, f2) {
			t.Errorf("invalid primary pointer mapping")
		}
	}
}

func TestPrimaryNilPointerStruct(t *testing.T) {
	p := primaryPtr{}

	var p2 primaryPtr
	bson.StructToBson(p).Struct(&p2)

	if p != p2 {
		t.Errorf("invalid primary nil pointer mapping")
	}
}

func TestPrimaryArraySliceStruct(t *testing.T) {
	type as struct {
		Array [3]int
		Slice []string
	}

	a := as{
		Array: [3]int{math.MinInt32, 0, math.MaxInt32},
		Slice: []string{"hello", "world"},
	}

	b := bson.StructToBson(a)

	var a2 as
	b.Struct(&a2)

	if a.Array != a2.Array {
		t.Errorf("invalid primary array mapping")
	}

	if len(a.Slice) != len(a2.Slice) {
		t.Errorf("invalid primary slice mapping")
	}

	for i := 0; i < len(a.Slice); i++ {
		if a.Slice[i] != a2.Slice[i] {
			t.Errorf("invalid primary slice mapping")
		}
	}
}

func TestPrimaryArraySlicePointerStruct(t *testing.T) {
	type as struct {
		Array *[3]int
		Slice *[]string
	}

	Array := [3]int{math.MinInt32, 0, math.MaxInt32}
	Slice := []string{"hello", "world"}

	a := as{
		Array: &Array,
		Slice: &Slice,
	}

	b := bson.StructToBson(a)

	var a2 as
	b.Struct(&a2)

	if *a.Array != *a2.Array {
		t.Errorf("invalid primary array pointer mapping")
	}

	if len(*a.Slice) != len(*a2.Slice) {
		t.Errorf("invalid primary slice pointer mapping")
	}

	for i := 0; i < len(*a.Slice); i++ {
		if (*a.Slice)[i] != (*a2.Slice)[i] {
			t.Errorf("invalid primary slice pointer mapping")
		}
	}
}

func TestPrimaryArraySliceNilPointerStruct(t *testing.T) {
	type as struct {
		Array *[3]int
		Slice *[]string
	}

	a := as{}

	b := bson.StructToBson(a)

	var a2 as
	b.Struct(&a2)

	if a != a2 {
		t.Errorf("invalid primary array pointer mapping")
	}
}

func TestPrimaryMapStruct(t *testing.T) {
	type mst struct {
		Map map[string]interface{}
	}

	m := mst{mdata}

	b := bson.StructToBson(&m)

	var m2 mst
	b.Struct(&m2)

	if len(m.Map) != len(m2.Map) {
		t.Errorf("invalid primary map mapping")
	}

	for k, v1 := range m.Map {
		v2, exist := m2.Map[k]
		if !exist {
			t.Errorf("invalid primary map mapping")
		}

		if !valueEqual(v1, v2) {
			t.Errorf("invalid primary map mapping")
		}
	}
}

func TestPrimaryMapPointerStruct(t *testing.T) {
	type mst struct {
		Map *map[string]interface{}
	}

	m := mst{&mdata}

	b := bson.StructToBson(&m)

	var m2 mst
	b.Struct(&m2)

	if len(*m.Map) != len(*m2.Map) {
		t.Errorf("invalid primary map pointer mapping")
	}

	for k, v := range *m.Map {
		v2, exist := (*m2.Map)[k]
		if !exist {
			t.Errorf("invalid primary map pointer mapping")
		}

		if !valueEqual(v, v2) {
			t.Errorf("invalid primary map pointer mapping")
		}
	}
}

func TestPrimaryMapNilPointerStruct(t *testing.T) {
	type mst struct {
		Map *map[string]interface{}
	}

	m := mst{}

	b := bson.StructToBson(&m)

	var m2 mst
	b.Struct(&m2)

	if m != m2 {
		t.Errorf("invalid primary map nil pointer mapping")
	}
}

func TestPrimaryDocStruct(t *testing.T) {
	type dst struct {
		Doc bson.Doc
	}

	d := dst{ddata}

	b := bson.StructToBson(&d)

	var d2 dst
	b.Struct(&d2)

	if len(d.Doc) != len(d2.Doc) {
		t.Errorf("invalid primary doc mapping")
	}

	for i := 0; i < len(d.Doc); i++ {
		e := d.Doc[i]
		e2 := d2.Doc[i]

		if e.Name != e2.Name {
			t.Errorf("invalid primary doc mapping")
		}

		if !valueEqual(e.Value, e2.Value) {
			t.Errorf("invalid primary doc mapping")
		}
	}
}

func TestPrimaryDocPointerStruct(t *testing.T) {
	type dst struct {
		Doc *bson.Doc
	}

	d := dst{&ddata}

	b := bson.StructToBson(&d)

	var d2 dst
	b.Struct(&d2)

	if len(*d.Doc) != len(*d2.Doc) {
		t.Errorf("invalid primary doc pointer mapping")
	}

	for i := 0; i < len(*d.Doc); i++ {
		e := (*d.Doc)[i]
		e2 := (*d2.Doc)[i]

		if e.Name != e2.Name {
			t.Errorf("invalid primary doc pointer mapping")
		}

		if !valueEqual(e.Value, e2.Value) {
			t.Errorf("invalid primary doc pointer mapping")
		}
	}
}

func TestPrimaryDocNilPointerStruct(t *testing.T) {
	type dst struct {
		Doc *bson.Doc
	}

	d := dst{}

	b := bson.StructToBson(&d)

	var d2 dst
	b.Struct(&d2)

	if d != d2 {
		t.Errorf("invalid nil doc mapping")
	}
}

func TestInterfaceStruct(t *testing.T) {
	type primary struct {
		Bool    interface{}
		Int     interface{}
		Int8    interface{}
		Int16   interface{}
		Int32   interface{}
		Int64   interface{}
		Uint    interface{}
		Uint8   interface{}
		Uint16  interface{}
		Uint32  interface{}
		Uint64  interface{}
		Uintptr interface{}
		Float32 interface{}
		Float64 interface{}
		String  interface{}
	}

	s := primary{
		Bool:    true,
		Int:     math.MinInt32,
		Int8:    math.MinInt8,
		Int16:   math.MinInt16,
		Int32:   math.MinInt32,
		Int64:   math.MinInt64,
		Uint:    math.MaxUint32,
		Uint8:   math.MaxUint8,
		Uint16:  math.MaxUint16,
		Uint32:  math.MaxUint32,
		Uint64:  math.MaxInt64,
		Uintptr: math.MaxInt32,
		Float32: math.MaxFloat32,
		Float64: math.MaxFloat64,
		String:  "string",
	}

	b := bson.StructToBson(s)

	var s2 primary
	b.Struct(&s2)

	value := reflect.ValueOf(s)
	value2 := reflect.ValueOf(s2)
	typ := reflect.TypeOf(primary{})
	for i := 0; i < typ.NumField(); i++ {
		f := value.Field(i).Interface()
		f2 := value2.Field(i).Interface()

		if !valueEqual(f, f2) {
			t.Errorf("invalid interface mapping")
		}
	}
}

func TestInterfacePointerStruct(t *testing.T) {
	type ptr struct {
		Bool    *interface{}
		Int     *interface{}
		Int8    *interface{}
		Int16   *interface{}
		Int32   *interface{}
		Int64   *interface{}
		Uint    *interface{}
		Uint8   *interface{}
		Uint16  *interface{}
		Uint32  *interface{}
		Uint64  *interface{}
		Uintptr *interface{}
		Float32 *interface{}
		Float64 *interface{}
		String  *interface{}
	}

	Bool := interface{}(true)
	Int := interface{}(math.MinInt32)
	Int8 := interface{}(math.MinInt8)
	Int16 := interface{}(math.MinInt16)
	Int32 := interface{}(math.MinInt32)
	Int64 := interface{}(math.MinInt64)
	Uint := interface{}(math.MaxUint32)
	Uint8 := interface{}(math.MaxUint8)
	Uint16 := interface{}(math.MaxUint16)
	Uint32 := interface{}(math.MaxUint32)
	Uint64 := interface{}(math.MaxInt64)
	Uintptr := interface{}(math.MaxInt32)
	Float32 := interface{}(math.MaxFloat32)
	Float64 := interface{}(math.MaxFloat64)
	String := interface{}("string")

	p := ptr{
		Bool:    &Bool,
		Int:     &Int,
		Int8:    &Int8,
		Int16:   &Int16,
		Int32:   &Int32,
		Int64:   &Int64,
		Uint:    &Uint,
		Uint8:   &Uint8,
		Uint16:  &Uint16,
		Uint32:  &Uint32,
		Uint64:  &Uint64,
		Uintptr: &Uintptr,
		Float32: &Float32,
		Float64: &Float64,
		String:  &String,
	}

	var p2 ptr
	b := bson.StructToBson(p)
	b.Struct(&p2)

	value := reflect.ValueOf(p)
	value2 := reflect.ValueOf(p2)
	typ := reflect.TypeOf(ptr{})
	for i := 0; i < typ.NumField(); i++ {
		f := value.Field(i).Interface()
		f2 := value2.Field(i).Interface()

		if !valueEqual(f, f2) {
			t.Errorf("invalid interface pointer mapping")
		}
	}
}

func TestInterfaceNilPointerStruct(t *testing.T) {
	type ptr struct {
		Bool    *interface{}
		Int     *interface{}
		Int8    *interface{}
		Int16   *interface{}
		Int32   *interface{}
		Int64   *interface{}
		Uint    *interface{}
		Uint8   *interface{}
		Uint16  *interface{}
		Uint32  *interface{}
		Uint64  *interface{}
		Uintptr *interface{}
		Float32 *interface{}
		Float64 *interface{}
		String  *interface{}
	}

	p := ptr{}

	var p2 ptr
	b := bson.StructToBson(p)
	b.Struct(&p2)

	if p != p2 {
		t.Errorf("invalid interface nil pointer mapping")
	}
}

func TestInterfaceStructStruct(t *testing.T) {
	type st2 struct {
		S primary
	}

	type st struct {
		S interface{}
	}

	pri := pridata
	s := st{pri}

	b := bson.StructToBson(s)

	var s2 st2
	b.Struct(&s2)

	if s2.S != pri {
		t.Errorf("invalid interface struct struct mapping")
	}
}

func TestStructStruct(t *testing.T) {
	type st struct {
		Primary primary
	}

	s := st{
		Primary: pridata,
	}

	b := bson.StructToBson(s)

	var s2 st
	b.Struct(&s2)

	if s.Primary != s2.Primary {
		t.Errorf("invalid struct struct mapping")
	}
}

func TestStructPointerStruct(t *testing.T) {
	type st struct {
		Primary *primary
	}

	s := st{
		Primary: &pridata,
	}

	b := bson.StructToBson(s)

	var s2 st
	b.Struct(&s2)

	if *s.Primary != *s2.Primary {
		t.Errorf("invalid struct pointer struct mapping")
	}
}

func TestStructNilPointerStruct(t *testing.T) {
	type st struct {
		Primary *primary
	}

	s := st{}

	b := bson.StructToBson(s)

	var s2 st
	b.Struct(&s2)

	if s != s2 {
		t.Errorf("invalid struct nil pointer struct mapping")
	}
}

func TestStructArraySliceStruct(t *testing.T) {
	type as struct {
		Array [3]int
		Slice []string
	}

	type st struct {
		As as
	}

	s := st{
		as{
			Array: [3]int{math.MinInt32, 0, math.MaxInt32},
			Slice: []string{"hello", "world"},
		},
	}

	b := bson.StructToBson(s)

	var s2 st
	b.Struct(&s2)

	if s.As.Array != s2.As.Array {
		t.Errorf("invalid struct array mapping")
	}

	if len(s.As.Slice) != len(s2.As.Slice) {
		t.Errorf("invalid struct slice mapping")
	}

	for i := 0; i < len(s.As.Slice); i++ {
		if s.As.Slice[i] != s2.As.Slice[i] {
			t.Errorf("invalid struct slice mapping")
		}
	}
}

func TestStructArraySlicePointerStruct(t *testing.T) {
	type as struct {
		Array *[3]int
		Slice *[]string
	}

	type st struct {
		As *as
	}

	s := st{
		&as{
			Array: &[3]int{math.MinInt32, 0, math.MaxInt32},
			Slice: &[]string{"hello", "world"},
		},
	}

	b := bson.StructToBson(s)

	var s2 st
	b.Struct(&s2)

	if (*(s.As.Array)) != (*(s2.As.Array)) {
		t.Errorf("invalid struct array pointer mapping")
	}

	if len(*(s.As.Slice)) != len(*(s2.As.Slice)) {
		t.Errorf("invalid struct slice pointer mapping")
	}

	for i := 0; i < len(*(s.As.Slice)); i++ {
		if (*(s.As.Slice))[i] != (*(s2.As.Slice))[i] {
			t.Errorf("invalid struct slice pointer mapping")
		}
	}
}

func TestStructArraySliceNilPointerStruct(t *testing.T) {
	type as struct {
		Array [3]int
		Slice []string
	}

	type st struct {
		As *as
	}

	s := st{}

	b := bson.StructToBson(s)

	var s2 st
	b.Struct(&s2)

	if s != s2 {
		t.Errorf("invalid struct array and slice pointer mapping")
	}
}

func TestStructMapStruct(t *testing.T) {
	type mst struct {
		Map map[string]interface{}
	}

	type st struct {
		M mst
	}

	m := st{
		mst{mdata},
	}

	b := bson.StructToBson(&m)

	var m2 st
	b.Struct(&m2)

	if len(m.M.Map) != len(m2.M.Map) {
		t.Errorf("invalid struct map mapping")
	}

	for k, v := range m.M.Map {
		v2, exist := (m2.M.Map)[k]
		if !exist {
			t.Errorf("invalid struct map mapping")
		}

		if !valueEqual(v, v2) {
			t.Errorf("invalid struct map mapping")
		}
	}
}

func TestStructMapPointerStruct(t *testing.T) {
	type mst struct {
		Map *map[string]interface{}
	}

	type st struct {
		M *mst
	}

	m := st{
		&mst{&mdata},
	}

	b := bson.StructToBson(&m)

	var m2 st
	b.Struct(&m2)

	if len(*(m.M.Map)) != len(*(m2.M.Map)) {
		t.Errorf("invalid struct map pointer mapping")
	}

	for k, v := range *(m.M.Map) {
		v2, exist := (*(m2.M.Map))[k]
		if !exist {
			t.Errorf("invalid struct map pointer mapping")
		}

		if !valueEqual(v, v2) {
			t.Errorf("invalid struct map pointer mapping")
		}
	}
}

func TestStructDocStruct(t *testing.T) {
	type dst struct {
		Doc bson.Doc
	}

	type st struct {
		D dst
	}

	d := st{
		dst{ddata},
	}

	b := bson.StructToBson(&d)

	var d2 st
	b.Struct(&d2)

	if len(d.D.Doc) != len(d2.D.Doc) {
		t.Errorf("invalid struct doc mapping")
	}

	for i := 0; i < len(d.D.Doc); i++ {
		e := d.D.Doc[i]
		e2 := d2.D.Doc[i]

		if e.Name != e2.Name {
			t.Errorf("invalid struct doc mapping")
		}

		if !valueEqual(e.Value, e2.Value) {
			t.Errorf("invalid struct doc mapping")
		}
	}
}

func TestStructDocPointerStruct(t *testing.T) {
	type dst struct {
		Doc *bson.Doc
	}

	type st struct {
		D *dst
	}

	d := st{
		&dst{&ddata},
	}

	b := bson.StructToBson(&d)

	var d2 st
	b.Struct(&d2)

	if len(*(d.D.Doc)) != len(*(d2.D.Doc)) {
		t.Errorf("invalid struct doc pointer mapping")
	}

	for i := 0; i < len(*(d.D.Doc)); i++ {
		e := (*(d.D.Doc))[i]
		e2 := (*(d2.D.Doc))[i]

		if e.Name != e2.Name {
			t.Errorf("invalid struct doc pointer mapping")
		}

		if !valueEqual(e.Value, e2.Value) {
			t.Errorf("invalid struct doc pointer mapping")
		}
	}
}

func TestArraySliceStructStruct(t *testing.T) {
	type st struct {
		Array [3]primary
		Slice []primary
	}

	s := st{
		Array: [3]primary{pridata, pridata, pridata},
		Slice: []primary{pridata, pridata, pridata},
	}

	b := bson.StructToBson(&s)

	var s2 st
	b.Struct(&s2)

	if s2.Array != s.Array {
		t.Errorf("invalid array struct mapping")
	}

	if len(s2.Slice) != len(s.Slice) {
		t.Fatalf("invalid slice struct mapping")
	}

	for i := 0; i < len(s.Slice); i++ {
		p := s.Slice[i]
		p2 := s2.Slice[i]

		if p != p2 {
			t.Errorf("invalid slice struct mapping")
		}
	}
}

func TestArraySliceDocStruct(t *testing.T) {
	type st struct {
		Array [3]bson.Doc
		Slice []bson.Doc
	}

	s := st{
		Array: [3]bson.Doc{ddata, ddata, ddata},
		Slice: []bson.Doc{ddata, ddata, ddata},
	}

	b := bson.StructToBson(s)

	var s2 st
	b.Struct(&s2)

	if len(s.Array) != len(s2.Array) {
		t.Errorf("invalid array doc mapping")
	}

	for i := 0; i < len(s.Array); i++ {
		a1 := s.Array[i]
		a2 := s2.Array[i]
		if !docValueEqual(a1, a2) {
			t.Errorf("invalid array doc mapping")
		}
	}

	if len(s.Slice) != len(s2.Slice) {
		t.Errorf("invalid slice doc mapping")
	}

	for i := 0; i < len(s.Slice); i++ {
		a1 := s.Slice[i]
		a2 := s2.Slice[i]
		if !docValueEqual(a1, a2) {
			t.Errorf("invalid slice doc mapping")
		}
	}
}

func BenchmarkSdbBsonMap(t *testing.B) {
	for i := 0; i < t.N; i++ {
		m := mdata2
		b := m.Bson()
		_ = b

		//m2 := b.Map()
		//_ = m2
	}
}

func BenchmarkSdbBsonDoc(t *testing.B) {
	for i := 0; i < t.N; i++ {
		doc := ddata
		b := doc.Bson()
		_ = b

		//doc2 := doc.Bson().Doc()
		//if len(doc) != len(doc2) {
		//	t.Errorf("bson fields num: %d, doc fileds num: %d", len(doc2), len(doc))
		//}
	}
}

func BenchmarkSdbBsonStruct(t *testing.B) {
	for i := 0; i < t.N; i++ {
		s := pridata
		b := bson.StructToBson(s)

		//var s2 primary
		//b.Struct(&s2)
		_ = b
	}
}

func BenchmarkSdbBsonAppendXXX(t *testing.B) {
	for i := 0; i < t.N; i++ {
		b := bson.NewBsonBuilder()
		b.AppendBool("bool", pridata.Bool)
		b.AppendInt64("int", int64(pridata.Int))
		b.AppendInt32("int8", int32(pridata.Int8))
		b.AppendInt32("int16", int32(pridata.Int16))
		b.AppendInt32("int32", int32(pridata.Int32))
		b.AppendInt64("int64", int64(pridata.Int64))
		b.AppendInt64("uint", int64(pridata.Uint))
		b.AppendInt32("uint8", int32(pridata.Uint8))
		b.AppendInt32("uint16", int32(pridata.Uint16))
		b.AppendInt32("uint32", int32(pridata.Uint32))
		b.AppendInt64("uint64", int64(pridata.Uint64))
		b.AppendInt64("uintprt", int64(pridata.Uintptr))
		b.AppendFloat64("float32", float64(pridata.Float32))
		b.AppendFloat64("float64", pridata.Float64)
		b.AppendString("string", pridata.String)
		b.Finish()

		//var s2 primary
		//b.Doc()
	}
}

func BenchmarkSdbBsonAppend(t *testing.B) {
	for i := 0; i < t.N; i++ {
		b := bson.NewBsonBuilder()
		for _, item := range ddata {
			b.Append(item.Name, item.Value)
		}
		b.Finish()

		//var s2 primary
		//b.Struct(&s2)
		//b.Doc()
	}
}
