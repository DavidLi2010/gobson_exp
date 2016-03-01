// Copyright 2015 David Li
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

import "fmt"

type Bson struct {
	raw      []byte
	order    ByteOrder
	offset   int // start position in raw
	inChild  bool
	finished bool
}

const initialBufferSize = 64
const eod = byte(0x00) // end of doc

func (bson *Bson) byteOrder() ByteOrder {
	return bson.order
}

func (bson *Bson) appendCString(value string) {
	const eos byte = 0x00 // end of cstring
	bson.raw = append(bson.raw, []byte(value)...)
	bson.raw = append(bson.raw, eos)
}

func (bson *Bson) reserveInt32() (pos int) {
	pos = len(bson.raw)
	bson.raw = append(bson.raw, 0, 0, 0, 0)
	return pos
}

func (bson *Bson) setInt32(pos int, value int32) {
	bson.order.SetInt32(bson.raw, pos, value)
}

func NewBson() *Bson {
	return NewBsonWithByteOrder(GetByteOrder())
}

func NewBsonWithByteOrder(order ByteOrder) *Bson {
	bson := &Bson{raw: make([]byte, 0, initialBufferSize), order: order}
	bson.reserveInt32()
	return bson
}

func (bson *Bson) checkBeforeAppend() {
	if bson.finished {
		panic("the bson is finished")
	}

	if bson.inChild {
		panic("in child bson")
	}
}

func (bson *Bson) Finish() {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, eod)
	bson.setInt32(bson.offset, int32(len(bson.raw)-bson.offset))
	bson.finished = true
}

func (bson *Bson) Raw() []byte {
	return bson.raw[bson.offset:]
}

func (bson *Bson) AppendDouble(name string, value float64) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeDouble))
	bson.appendCString(name)
	bson.raw = bson.order.AppendFloat64(bson.raw, value)
}

func (bson *Bson) AppendString(name string, value string) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeString))
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt32(bson.raw, int32(len(value)+1))
	bson.appendCString(value)
}

func (bson *Bson) AppendBson(name string, value *Bson) {
	bson.checkBeforeAppend()
	if bson.order != value.order {
		panic("the byte order is different")
	}
	bson.raw = append(bson.raw, byte(BsonTypeDoc))
	bson.appendCString(name)
	bson.raw = append(bson.raw, value.Raw()...)
}

func (bson *Bson) AppendBsonStart(name string) (child *Bson) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeDoc))
	bson.appendCString(name)
	child = &Bson{raw: bson.raw, order: bson.order, offset: len(bson.raw)}
	bson.inChild = true
	return child
}

func (bson *Bson) AppendBsonEnd() {
	if !bson.inChild {
		panic("not in child bson")
	}
	if bson.finished {
		panic("the bson is finished")
	}
	if bson.raw[len(bson.raw)-1] != eod {
		panic("the child bson is not finished")
	}
	bson.inChild = false
}

func (bson *Bson) AppendArray(name string, value *BsonArray) {
	bson.checkBeforeAppend()
	if bson.order != value.bson.order {
		panic("the byte order is different")
	}
	bson.raw = append(bson.raw, byte(BsonTypeArray))
	bson.appendCString(name)
	bson.raw = append(bson.raw, value.Raw()...)
}

func (bson *Bson) AppendArrayStart(name string) (child *BsonArray) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeArray))
	bson.appendCString(name)
	child = &BsonArray{bson: Bson{raw: bson.raw, order: bson.order, offset: len(bson.raw)}}
	bson.inChild = true
	return child
}

func (bson *Bson) AppendArrayEnd() {
	if !bson.inChild {
		panic("not in child array")
	}
	if bson.finished {
		panic("the bson is finished")
	}
	if bson.raw[len(bson.raw)-1] != eod {
		panic("the child array is not finished")
	}
	bson.inChild = false
}

func (bson *Bson) AppendBinary(name string, value Binary) {
	bson.checkBeforeAppend()
	if value.Data == nil {
		panic("binary is null")
	}
	bson.raw = append(bson.raw, byte(BsonTypeBinary))
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt32(bson.raw, int32(len(value.Data)))
	bson.raw = append(bson.raw, byte(value.Subtype))
	bson.raw = append(bson.raw, value.Data...)
}

func (bson *Bson) AppendObjectId(name string, value ObjectId) {
	bson.checkBeforeAppend()
	if !value.Valid() {
		panic(fmt.Sprintf("invalid ObjectId: %s", value))
	}
	bson.raw = append(bson.raw, byte(BsonTypeObjectId))
	bson.appendCString(name)
	bson.raw = append(bson.raw, []byte(value)...)
}

func (bson *Bson) AppendBool(name string, value bool) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeBool))
	bson.appendCString(name)
	if value {
		bson.raw = append(bson.raw, byte(1))
	} else {
		bson.raw = append(bson.raw, byte(0))
	}
}

func (bson *Bson) AppendDate(name string, value Date) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeDate))
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt64(bson.raw, int64(value))
}

func (bson *Bson) AppendNull(name string) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeNull))
	bson.appendCString(name)
}

func (bson *Bson) AppendRegex(name string, pattern string, options string) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeRegEx))
	bson.appendCString(name)
	bson.appendCString(pattern)
	bson.appendCString(options)
}

func (bson *Bson) AppendInt32(name string, value int32) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeInt32))
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt32(bson.raw, value)
}

func (bson *Bson) AppendTimestamp(name string, value Timestamp) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeTimestamp))
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt32(bson.raw, value.Increment)
	bson.raw = bson.order.AppendInt32(bson.raw, value.Second)
}

func (bson *Bson) AppendInt64(name string, value int64) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeInt64))
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt64(bson.raw, value)
}

func (bson *Bson) AppendMinKey(name string) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeMinKey))
	bson.appendCString(name)
}

func (bson *Bson) AppendMaxKey(name string) {
	bson.checkBeforeAppend()
	bson.raw = append(bson.raw, byte(BsonTypeMaxKey))
	bson.appendCString(name)
}
