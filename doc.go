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

const (
	// bson type
	bsonTypeEOD        = byte(0x00) // end of doc
	bsonTypeDouble     = byte(0x01)
	bsonTypeString     = byte(0x02)
	bsonTypeDoc        = byte(0x03)
	bsonTypeArray      = byte(0x04)
	bsonTypeBinary     = byte(0x05)
	bsonTypeUndefined  = byte(0x06) // deprecated
	bsonTypeObjectId   = byte(0x07)
	bsonTypeBool       = byte(0x08)
	bsonTypeDate       = byte(0x09)
	bsonTypeNull       = byte(0x0A)
	bsonTypeRegex      = byte(0x0B)
	bsonTypeDBPointer  = byte(0x0C) // deprecated
	bsonTypeCode       = byte(0x0D)
	bsonTypeSymbol     = byte(0x0E) // deprecated
	bsonTypeCodeWScope = byte(0x0F)
	bsonTypeInt32      = byte(0x10)
	bsonTypeTimestamp  = byte(0x11)
	bsonTypeInt64      = byte(0x12)
	bsonTypeMaxKey     = byte(0x7F)
	bsonTypeMinKey     = byte(0xFF)

	binaryTypeGeneral          = byte(0x00)
	binaryTypeFunction         = byte(0x01)
	binaryTypeBinaryDeprecated = byte(0x02)
	binaryTypeUUIDDeprecated   = byte(0x03)
	binaryTypeUUID             = byte(0x00)
	binaryTypeMD5              = byte(0x05)
	binaryTypeUser             = byte(0x80)

	// end of cstring
	eos = byte(0x00)
)

type Bson struct {
	raw   []byte
	order ByteOrder
}

type Date int64

type RegEx struct {
	pattern string
	options string
}

type Timestamp struct {
	Second    int32
	Increment int32
}

type Binary struct {
	Subtype byte
	Data    []byte
}

func (bson *Bson) byteOrder() ByteOrder {
	return bson.order
}

func (bson *Bson) appendCString(value string) {
	bson.raw = append(bson.raw, []byte(value)...)
	bson.raw = append(bson.raw, eos)
}

func (bson *Bson) appendDouble(name string, value float64) {
	bson.raw = append(bson.raw, bsonTypeDouble)
	bson.appendCString(name)
	bson.raw = bson.order.AppendFloat64(bson.raw, value)
}

func (bson *Bson) appendString(name string, value string) {
	bson.raw = append(bson.raw, bsonTypeString)
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt32(bson.raw, int32(len(value)+1))
	bson.appendCString(value)
}

func (bson *Bson) appendDoc(name string, value *Bson) {
	if bson.order != value.order {
		panic("the byte order is different")
	}
	bson.raw = append(bson.raw, bsonTypeDoc)
	bson.appendCString(name)
	bson.raw = append(bson.raw, value.raw...)
}

func (bson *Bson) appendArray(name string, value []Bson) {

}

func (bson *Bson) appendBinary(name string, value Binary) {
	bson.raw = append(bson.raw, bsonTypeBinary)
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt32(len(value.Data))
	bson.raw = append(bson.raw, value.Subtype)
	bson.raw = append(bson.raw, value.Data...)
}

func (bson *Bson) appendObjectId(name string, value ObjectId) {
	if !value.Valid() {
		panic(fmt.Sprintf("invalid ObjectId: %s", value))
	}
	bson.raw = append(bson.raw, bsonTypeObjectId)
	bson.appendCString(name)
	bson.raw = append(bson.raw, []byte(value)...)
}

func (bson *Bson) appendBool(name string, value bool) {
	bson.raw = append(bson.raw, bsonTypeBool)
	bson.appendCString(name)
	if value {
		bson.raw = append(bson.raw, byte(1))
	} else {
		bson.raw = append(bson.raw, byte(0))
	}
}

func (bson *Bson) appendDate(name string, value Date) {
	bson.raw = append(bson.raw, bsonTypeDate)
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt64(int64(value))
}

func (bson *Bson) appendNull(name string) {
	bson.raw = append(bson.raw, bsonTypeNull)
	bson.appendCString(name)
}

func (bson *Bson) appendRegex(name string, pattern string, options string) {
	bson.raw = append(bson.raw, bsonTypeRegex)
	bson.appendCString(name)
	bson.appendCString(pattern)
	bson.appendCString(options)
}

func (bson *Bson) appendInt32(name string, value int32) {
	bson.raw = append(bson.raw, bsonTypeInt32)
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt32(bson.raw, value)
}

func (bson *Bson) appendTimestamp(name string, value Timestamp) {
	bson.raw = append(bson.raw, bsonTypeTimestamp)
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt32(bson.raw, value.Increment)
	bson.raw = bson.order.AppendInt32(bson.raw, value.Second)
}

func (bson *Bson) appendInt64(name string, value int64) {
	bson.raw = append(bson.raw, bsonTypeTimestamp)
	bson.appendCString(name)
	bson.raw = bson.order.AppendInt64(bson.raw, value)
}

func (bson *Bson) appendMinKey(name string) {
	bson.raw = append(bson.raw, bsonTypeMinKey)
	bson.appendCString(name)
}

func (bson *Bson) appendMaxKey(name string) {
	bson.raw = append(bson.raw, bsonTypeMaxKey)
	bson.appendCString(name)
}
