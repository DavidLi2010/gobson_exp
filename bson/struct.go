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

package bson

import (
	"fmt"
	"reflect"
)

type structInfo struct {
	Fields []fieldInfo
}

type fieldInfo struct {
	Name  string
	Index int
}

func structToBsonBuilder(s reflect.Value, b *BsonBuilder) {
	t := s.Type()
	n := s.NumField()
	for i := 0; i < n; i++ {
		field := t.Field(i)
		// private field
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}
		b.Append(field.Name, s.Field(i).Interface())
	}
}

func StructToBson(s interface{}) *Bson {
	v := reflect.ValueOf(s)
	switch v.Kind() {
	case reflect.Struct:
	case reflect.Ptr:
		if v.Elem().Kind() != reflect.Struct {
			panic("s must be struct or struct pointer")
		}
		v = v.Elem()
	default:
		panic("s must be struct or struct pointer")
	}

	b := NewBsonBuilder()
	structToBsonBuilder(v, b)
	b.Finish()
	return b.Bson()
}

func mapToStruct(s reflect.Value, m Map) {
	t := s.Type()
	n := s.NumField()
	for i := 0; i < n; i++ {
		field := t.Field(i)
		// private field
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		v, exist := m[field.Name]
		if exist {
			setFieldValue(s.Field(i), v)
		}
	}
}

func docToStruct(s reflect.Value, d Doc) {
	t := s.Type()
	n := s.NumField()

	m := map[string]int{}
	for idx, v := range d {
		m[v.Name] = idx
	}

	for i := 0; i < n; i++ {
		field := t.Field(i)
		// private field
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		idx, exist := m[field.Name]
		if exist {
			setFieldValue(s.Field(i), d[idx].Value)
		}
	}
}

func tryConvert(in reflect.Value, t reflect.Type) (out reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	return in.Convert(t), nil
}

func setFieldValue(f reflect.Value, v interface{}) {
	value := reflect.ValueOf(v)
	switch f.Kind() {
	case reflect.Interface:
		f.Set(value)
	case reflect.String:
		switch value.Kind() {
		case reflect.String:
			f.SetString(value.String())
		case reflect.Slice:
			if b, ok := v.([]byte); ok {
				f.SetString(string(b))
			}
		}
	case reflect.Slice:
		switch value.Kind() {
		case reflect.Slice, reflect.Array:
			if f.Type().Elem() == value.Type().Elem() {
				if f.IsNil() {
					f.Set(reflect.MakeSlice(f.Type(), value.Len(), value.Len()))
				}
				reflect.Copy(f, value)
			} else {
				ft := f.Type().Elem()
				for i := 0; i < value.Len(); i++ {
					ev := value.Index(i)
					if ev.Kind() == reflect.Interface {
						ev = ev.Elem()
					}

					if ev.Type() == ft {
						f.Set(reflect.Append(f, ev))
					} else if ft.Kind() == reflect.Struct {
						et := ev.Type()
						if et == reflect.TypeOf(Doc{}) {
							fv := reflect.New(ft).Elem()
							docToStruct(fv, ev.Interface().(Doc))
							f.Set(reflect.Append(f, fv))
						} else if et == reflect.TypeOf(Map{}) {
							fv := reflect.New(ft).Elem()
							mapToStruct(fv, ev.Interface().(Map))
							f.Set(reflect.Append(f, fv))
						}
					} else {
						if cv, err := tryConvert(ev, ft); err == nil {
							f.Set(reflect.Append(f, cv))
						}
					}
				}
			}
		}
	case reflect.Array:
		switch value.Kind() {
		case reflect.Slice, reflect.Array:
			if f.Type().Elem() == value.Type().Elem() {
				reflect.Copy(f, value)
			} else {
				ft := f.Type().Elem()
				for i := 0; i < f.Len(); i++ {
					fev := f.Index(i)
					ev := value.Index(i)
					if ev.Kind() == reflect.Interface {
						ev = ev.Elem()
					}

					if ev.Type() == ft {
						fev.Set(ev)
					} else if ft.Kind() == reflect.Struct {
						et := ev.Type()
						if et == reflect.TypeOf(Doc{}) {
							fv := reflect.New(ft).Elem()
							docToStruct(fv, ev.Interface().(Doc))
							fev.Set(fv)
						} else if et == reflect.TypeOf(Map{}) {
							fv := reflect.New(ft).Elem()
							mapToStruct(fv, ev.Interface().(Map))
							fev.Set(fv)
						}
					} else {
						if cv, err := tryConvert(ev, f.Type().Elem()); err == nil {
							fev.Set(cv)
						}
					}
				}
			}

		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetInt(value.Int())
		case reflect.Float32, reflect.Float64:
			f.SetInt(int64(value.Float()))
		case reflect.Bool:
			if value.Bool() {
				f.SetInt(1)
			} else {
				f.SetInt(0)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			panic("no uint types in Bson")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetUint(uint64(value.Int()))
		case reflect.Float32, reflect.Float64:
			f.SetUint(uint64(value.Float()))
		case reflect.Bool:
			if value.Bool() {
				f.SetUint(1)
			} else {
				f.SetUint(0)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			panic("no uint types in Bson")
		}
	case reflect.Float32, reflect.Float64:
		switch value.Kind() {
		case reflect.Float32, reflect.Float64:
			f.SetFloat(value.Float())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetFloat(float64(value.Int()))
		case reflect.Bool:
			if value.Bool() {
				f.SetFloat(1)
			} else {
				f.SetFloat(0)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			panic("no uint types in Bson")
		}
	case reflect.Bool:
		switch value.Kind() {
		case reflect.Bool:
			f.SetBool(value.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			f.SetBool(value.Int() != 0)
		case reflect.Float32, reflect.Float64:
			f.SetBool(value.Float() != 0)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			panic("no uint types in Bson")
		}
	case reflect.Ptr:
		if v != nil && f.IsNil() {
			f.Set(reflect.New(f.Type().Elem()))
		}

		setFieldValue(f.Elem(), v)
	case reflect.Struct:
		switch value.Kind() {
		case reflect.Map:
			if m, ok := v.(Map); ok {
				mapToStruct(f, m)
			} /*else if d, ok := v.(Doc); ok {
				docToStruct(f, d)
			}*/
		case reflect.Slice:
			if d, ok := v.(Doc); ok {
				docToStruct(f, d)
			}
		}
	case reflect.Map:
		if f.Type().Key().Kind() != reflect.String {
			panic("the key of map must be string")
		}
		switch value.Kind() {
		case reflect.Map:
			if f.IsNil() {
				f.Set(reflect.MakeMap(f.Type()))
			}
			if f.Type().Elem() == value.Type().Elem() {
				for _, k := range value.MapKeys() {
					f.SetMapIndex(k, value.MapIndex(k))
				}
			} else {
				for _, k := range value.MapKeys() {
					mv := value.MapIndex(k)
					if mv.Kind() == reflect.Interface {
						mv = mv.Elem()
					}
					if nv, err := tryConvert(mv, f.Type().Elem()); err == nil {
						f.SetMapIndex(k, nv)
					}
				}
			}
		case reflect.Slice:
			if value.Type() != reflect.TypeOf(Doc{}) {
				return
			}
			d := v.(Doc)
			if f.IsNil() {
				f.Set(reflect.MakeMap(f.Type()))
			}
			if f.Type().Elem() == value.Type().Elem() {
				for _, ev := range d {
					f.SetMapIndex(reflect.ValueOf(ev.Name), reflect.ValueOf(ev.Value))
				}
			} else {
				for _, ev := range d {
					mv := reflect.ValueOf(ev.Value)
					if mv.Kind() == reflect.Interface {
						mv = mv.Elem()
					}
					if nv, err := tryConvert(mv, f.Type().Elem()); err == nil {
						f.SetMapIndex(reflect.ValueOf(ev.Name), nv)
					}
				}
			}
		}
	}
}
