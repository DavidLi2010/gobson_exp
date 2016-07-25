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
	"math"
	"reflect"
)

type BsonBuilder struct {
	raw      []byte
	offset   int
	child    *BsonBuilder
	parent   *BsonBuilder
	inChild  bool
	finished bool
}

func NewBsonBuilder() *BsonBuilder {
	b := &BsonBuilder{raw: make([]byte, 0, initialBufferSize)}
	b.reserveLength()
	return b
}

func (b *BsonBuilder) reserveLength() {
	b.raw = append(b.raw, 0, 0, 0, 0)
}

func (b *BsonBuilder) setLength(v int32) {
	b.raw[b.offset+0] = byte(v)
	b.raw[b.offset+1] = byte(v >> 8)
	b.raw[b.offset+2] = byte(v >> 16)
	b.raw[b.offset+3] = byte(v >> 24)
}

func (b *BsonBuilder) appendType(t BsonType) {
	b.raw = append(b.raw, byte(t))
}

func (b *BsonBuilder) appendCString(v string) {
	const eos byte = 0x00 // end of cstring
	b.raw = append(b.raw, []byte(v)...)
	b.raw = append(b.raw, eos)
}

func (b *BsonBuilder) appendBytes(v ...byte) {
	b.raw = append(b.raw, v...)
}

func (b *BsonBuilder) appendInt32(v int32) {
	u := uint32(v)
	b.raw = append(b.raw, byte(u), byte(u>>8), byte(u>>16), byte(u>>24))
}

func (b *BsonBuilder) appendInt64(v int64) {
	u := uint64(v)
	b.raw = append(b.raw, byte(u), byte(u>>8), byte(u>>16), byte(u>>24),
		byte(u>>32), byte(u>>40), byte(u>>48), byte(u>>56))
}

func (b *BsonBuilder) appendFloat64(v float64) {
	u := int64(math.Float64bits(v))
	b.appendInt64(u)
}

func (b *BsonBuilder) checkBeforeAppend() {
	if b.finished {
		panic("the bson builder is finished")
	}

	if b.inChild {
		panic("in child bson builder")
	}
}

func (b *BsonBuilder) Finish() *BsonBuilder {
	b.checkBeforeAppend()
	b.raw = append(b.raw, eod)
	b.setLength(int32(len(b.raw) - b.offset))
	b.finished = true
	return b
}

func (b *BsonBuilder) Raw() []byte {
	return b.raw[b.offset:]
}

func (b *BsonBuilder) Bson() *Bson {
	return &Bson{raw: b.Raw()}
}

func (b *BsonBuilder) AppendFloat64(name string, value float64) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeFloat64)
	b.appendCString(name)
	b.appendFloat64(value)
	return b
}

func (b *BsonBuilder) AppendString(name string, value string) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeString)
	b.appendCString(name)
	b.appendInt32(int32(len(value) + 1))
	b.appendCString(value)
	return b
}

func (b *BsonBuilder) AppendBson(name string, value *Bson) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeBson)
	b.appendCString(name)
	b.appendBytes(value.Raw()...)
	return b
}

func (parent *BsonBuilder) AppendBsonStart(name string) (child *BsonBuilder) {
	parent.checkBeforeAppend()
	parent.appendType(BsonTypeBson)
	parent.appendCString(name)
	child = &BsonBuilder{raw: parent.raw, offset: len(parent.raw)}
	child.reserveLength()
	parent.inChild = true
	parent.child = child
	child.parent = parent
	return child
}

func (child *BsonBuilder) AppendBsonEnd() (parent *BsonBuilder) {
	if child.parent == nil {
		panic("not in child bson builder")
	}
	if !child.finished {
		panic("the child bson builder is not finished")
	}
	if child.raw[len(child.raw)-1] != eod {
		panic("the child bson builder is not finished")
	}
	parent = child.parent
	parent.raw = child.raw
	child.parent = nil
	parent.child = nil
	parent.inChild = false
	return parent
}

func (b *BsonBuilder) AppendArray(name string, value *BsonArray) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeArray)
	b.appendCString(name)
	b.appendBytes(value.Raw()...)
	return b
}

func (parent *BsonBuilder) AppendArrayStart(name string) (child *BsonArrayBuilder) {
	parent.checkBeforeAppend()
	parent.appendType(BsonTypeArray)
	parent.appendCString(name)
	child = &BsonArrayBuilder{builder: BsonBuilder{raw: parent.raw, offset: len(parent.raw)}}
	child.builder.reserveLength()
	child.builder.parent = parent
	parent.inChild = true
	parent.child = &child.builder
	return child
}

func (b *BsonBuilder) AppendBinary(name string, value Binary) *BsonBuilder {
	b.checkBeforeAppend()
	if value.Data == nil {
		panic("binary is null")
	}
	b.appendType(BsonTypeBinary)
	b.appendCString(name)
	b.appendInt32(int32(len(value.Data)))
	b.appendBytes(byte(value.Subtype))
	b.appendBytes(value.Data...)
	return b
}

func (b *BsonBuilder) AppendObjectId(name string, value ObjectId) *BsonBuilder {
	b.checkBeforeAppend()
	if !value.IsValid() {
		panic(fmt.Sprintf("invalid ObjectId: %s", value))
	}
	b.appendType(BsonTypeObjectId)
	b.appendCString(name)
	b.appendBytes([]byte(value)...)
	return b
}

func (b *BsonBuilder) AppendBool(name string, value bool) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeBool)
	b.appendCString(name)
	if value {
		b.appendBytes(byte(1))
	} else {
		b.appendBytes(byte(0))
	}
	return b
}

func (b *BsonBuilder) AppendDate(name string, value Date) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeDate)
	b.appendCString(name)
	b.appendInt64(int64(value))
	return b
}

func (b *BsonBuilder) AppendNull(name string) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeNull)
	b.appendCString(name)
	return b
}

func (b *BsonBuilder) AppendRegex(name string, value RegEx) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeRegEx)
	b.appendCString(name)
	b.appendCString(value.Pattern)
	b.appendCString(value.Options)
	return b
}

func (b *BsonBuilder) AppendInt32(name string, value int32) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeInt32)
	b.appendCString(name)
	b.appendInt32(value)
	return b
}

func (b *BsonBuilder) AppendTimestamp(name string, value Timestamp) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeTimestamp)
	b.appendCString(name)
	b.appendInt32(value.Increment)
	b.appendInt32(value.Second)
	return b
}

func (b *BsonBuilder) AppendInt64(name string, value int64) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeInt64)
	b.appendCString(name)
	b.appendInt64(value)
	return b
}

func (b *BsonBuilder) AppendMinKey(name string) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeMinKey)
	b.appendCString(name)
	return b
}

func (b *BsonBuilder) AppendMaxKey(name string) *BsonBuilder {
	b.checkBeforeAppend()
	b.appendType(BsonTypeMaxKey)
	b.appendCString(name)
	return b
}

func (bson *BsonBuilder) Append(name string, value interface{}) {
	switch value.(type) {
	case float32:
		bson.AppendFloat64(name, float64(value.(float32)))
	case float64:
		bson.AppendFloat64(name, value.(float64))
	case int8:
		bson.AppendInt32(name, int32(value.(int8)))
	case int16:
		bson.AppendInt32(name, int32(value.(int16)))
	case int32:
		bson.AppendInt32(name, value.(int32))
	case int64:
		val := value.(int64)
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			bson.AppendInt32(name, int32(val))
		} else {
			bson.AppendInt64(name, int64(val))
		}
	case uint8:
		bson.AppendInt32(name, int32(value.(uint8)))
	case uint16:
		bson.AppendInt32(name, int32(value.(uint16)))
	case uint32:
		val := value.(uint32)
		if int32(val) < 0 {
			bson.AppendInt64(name, int64(val))
		} else {
			bson.AppendInt32(name, int32(val))
		}
	case uint64:
		val := int64(value.(uint64))
		if val < 0 {
			panic("bson has no uint64 type, and value is too large to fit correctly in an int64")
		}
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			bson.AppendInt32(name, int32(val))
		} else {
			bson.AppendInt64(name, int64(val))
		}
	case int:
		val := int64(value.(int))
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			bson.AppendInt32(name, int32(val))
		} else {
			bson.AppendInt64(name, int64(val))
		}
	case uint:
		val := int64(value.(uint))
		if val < 0 {
			panic("bson has no uint64 type, and value is too large to fit correctly in an int64")
		}
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			bson.AppendInt32(name, int32(val))
		} else {
			bson.AppendInt64(name, int64(val))
		}
	case uintptr:
		val := int64(value.(uintptr))
		if val < 0 {
			panic("bson has no uint64 type, and value is too large to fit correctly in an int64")
		}
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			bson.AppendInt32(name, int32(val))
		} else {
			bson.AppendInt64(name, int64(val))
		}
	case bool:
		bson.AppendBool(name, value.(bool))
	case string:
		bson.AppendString(name, value.(string))
	case nil:
		bson.AppendNull(name)
	case ObjectId:
		bson.AppendObjectId(name, value.(ObjectId))
	case Date:
		bson.AppendDate(name, value.(Date))
	case RegEx:
		bson.AppendRegex(name, value.(RegEx))
	case Timestamp:
		bson.AppendTimestamp(name, value.(Timestamp))
	case Binary:
		bson.AppendBinary(name, value.(Binary))
	case orderKey:
		val := value.(orderKey)
		if val == MaxKey {
			bson.AppendMaxKey(name)
		} else if val == MinKey {
			bson.AppendMinKey(name)
		} else {
			panic("invalid orderkey")
		}
	case Map:
		m := value.(Map)
		child := bson.AppendBsonStart(name)
		m.toBsonBuilder(child)
		child.Finish()
		bson.AppendBsonEnd()
	case Doc:
		d := value.(Doc)
		child := bson.AppendBsonStart(name)
		d.toBsonBuilder(child)
		child.Finish()
		child.AppendBsonEnd()
	default:
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Array, reflect.Slice:
			l := v.Len()
			child := bson.AppendArrayStart(name)
			for i := 0; i < l; i++ {
				child.Append(v.Index(i).Interface())
			}
			child.Finish()
			child.AppendArrayEnd()
			return
		case reflect.Map:
			child := bson.AppendBsonStart(name)
			for _, k := range v.MapKeys() {
				child.Append(k.String(), v.MapIndex(k).Interface())
			}
			child.Finish()
			child.AppendBsonEnd()
			return
		case reflect.Ptr:
			if v.IsNil() {
				bson.AppendNull(name)
			} else {
				bson.Append(name, v.Elem().Interface())
			}
			return
		case reflect.Struct:
			child := bson.AppendBsonStart(name)
			structToBsonBuilder(v, child)
			child.Finish()
			child.AppendBsonEnd()
			return
		}
		// Complex64, Complex128
		// Chan, Func
		// UnsafePointer
		panic(fmt.Errorf("can't append %s(%v) to bson", reflect.TypeOf(value).String(), value))
	}
}
