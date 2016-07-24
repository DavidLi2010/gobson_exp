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
	"bytes"
	"fmt"
	"reflect"
)

type Bson struct {
	raw []byte
}

const initialBufferSize = 64
const eod = byte(0x00) // end of doc

func NewBson(raw []byte) *Bson {
	return &Bson{raw:raw}
}

func (bson *Bson) Validate() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid bson: %v", r)
		}
	}()

	if bson.Length() != len(bson.Raw()) {
		return fmt.Errorf("invalid bson length")
	}

	it := bson.Iterator()
	for it.Next() {
		it.Name()
		it.Value()
	}

	return nil
}

func (bson *Bson) Raw() []byte {
	return bson.raw
}

func (bson *Bson) Iterator() *BsonIterator {
	return NewBsonIterator(bson)
}

func (bson *Bson) Length() int {
	return int(bytesToInt32(bson.raw))
}

func (bson *Bson) String() string {
	var err error
	buf := bytes.NewBufferString("{")
	it := bson.Iterator()
	for it.Next() {
		switch it.BsonType() {
		case BsonTypeFloat64:
			_, err = fmt.Fprintf(buf, `"%s":%v`, it.Name(), it.Float64())
		case BsonTypeString:
			_, err = fmt.Fprintf(buf, `"%s":"%s"`, it.Name(), it.UTF8String())
		case BsonTypeBson:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), it.Bson().String())
		case BsonTypeArray:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), it.BsonArray().String())
		case BsonTypeBinary:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), it.Binary().String())
		case BsonTypeObjectId:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), it.ObjectId().String())
		case BsonTypeBool:
			_, err = fmt.Fprintf(buf, `"%s":%v`, it.Name(), it.Bool())
		case BsonTypeDate:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), it.Date().String())
		case BsonTypeNull:
			_, err = fmt.Fprintf(buf, `"%s":null`, it.Name())
		case BsonTypeRegEx:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), it.RegEx().String())
		case BsonTypeInt32:
			_, err = fmt.Fprintf(buf, `"%s":%v`, it.Name(), it.Int32())
		case BsonTypeTimestamp:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), it.Timestamp().String())
		case BsonTypeInt64:
			_, err = fmt.Fprintf(buf, `"%s":%v`, it.Name(), it.Int64())
		case BsonTypeMaxKey:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), MaxKey.String())
		case BsonTypeMinKey:
			_, err = fmt.Fprintf(buf, `"%s":%s`, it.Name(), MinKey.String())
		case BsonTypeEOD:
			// END
		case BsonTypeUndefined: // deprecated
			fallthrough
		case BsonTypeDBPointer: // deprecated
			fallthrough
		case BsonTypeCode: // not support
			fallthrough
		case BsonTypeSymbol: // deprecated
			fallthrough
		case BsonTypeCodeWScope: // not support
			fallthrough
		default:
			panic(fmt.Errorf("invalid bson type: %v", it.BsonType()))
		}

		if err != nil {
			panic(fmt.Sprintf("failed to convert bson to string: %v", err))
		}

		if it.More() {
			_, err = buf.WriteString(", ")
			if err != nil {
				panic(fmt.Sprintf("failed to convert bson to string: %v", err))
			}
		}
	}

	_, err = buf.WriteString("}")
	if err != nil {
		panic(fmt.Sprintf("failed to convert bson to string: %v", err))
	}

	return buf.String()
}

func (bson *Bson) Map() Map {
	m := make(Map)

	it := bson.Iterator()
	for it.Next() {
		switch it.BsonType() {
		case BsonTypeBson:
			m[it.Name()] = it.Bson().Map()
		case BsonTypeArray:
			m[it.Name()] = it.BsonArray().MapSlice()
		default:
			m[it.Name()] = it.Value()
		}
	}

	return m
}

func (bson *Bson) Doc() Doc {
	d := []DocElement{}

	it := bson.Iterator()
	for it.Next() {
		var val interface{}
		switch it.BsonType() {
		case BsonTypeBson:
			val = it.Bson().Doc()
		case BsonTypeArray:
			val = it.BsonArray().DocSlice()
		default:
			val = it.Value()
		}
		d = append(d, DocElement{Name: it.Name(), Value: val})
	}

	return d
}

func (bson *Bson) Struct(s interface{}) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("s must be struct pointer")
	}
	docToStruct(v.Elem(), bson.Doc())
}
