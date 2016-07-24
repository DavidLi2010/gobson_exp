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
)

type BsonIterator struct {
	raw    []byte
	length int
	offset int

	// current field
	elementLen int
	keyLen     int
	value      []byte
}

func bytesToInt32(b []byte) int32 {
	if len(b) < 4 {
		panic("bytesToInt32: len([]byte) < 4")
	}

	return int32(uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24))
}

func bytesToInt64(b []byte) int64 {
	if len(b) < 8 {
		panic("bytesToInt64: len([]byte) < 8")
	}

	return int64(uint64(b[0]) | (uint64(b[1]) << 8) | (uint64(b[2]) << 16) | (uint64(b[3]) << 24) |
		(uint64(b[4]) << 32) | (uint64(b[5]) << 40) | (uint64(b[6]) << 48) | (uint64(b[7]) << 56))
}

func bytesToFloat64(b []byte) float64 {
	if len(b) < 8 {
		panic("bytesToFloat64: len([]byte) < 8")
	}

	return math.Float64frombits(uint64(bytesToInt64(b)))
}

func NewBsonIterator(bson *Bson) *BsonIterator {
	if bson == nil {
		panic("null bson")
	}
	if !bson.finished {
		panic("unfinished bson")
	}

	raw := bson.Raw()
	length := int(bytesToInt32(raw))

	return &BsonIterator{
		raw:    raw,
		length: length,
		offset: 4,
	}
}

func (it *BsonIterator) Reset() {
	it.offset = 4
	it.elementLen = 0
	it.keyLen = 0
	it.value = nil
}

// include '0x00'
func cstringLength(s []byte) int {
	len := len(s)
	if len == 0 {
		panic("no cstring")
	}

	cstringLen := 0
	for i := 0; i < len; i++ {
		if s[i] == 0x00 {
			cstringLen = i + 1
			break
		}
	}

	if cstringLen < 1 || s[cstringLen-1] != 0x00 {
		panic("invalid cstring")
	}

	return cstringLen
}

func (it *BsonIterator) More() bool {
	t := BsonType(it.raw[it.offset+it.elementLen])
	return t != BsonTypeEOD
}

func (it *BsonIterator) Next() bool {
	it.offset += it.elementLen
	t := BsonType(it.raw[it.offset])
	if t == BsonTypeEOD {
		it.elementLen = 0
		return false
	}

	fieldOffset := 1 // skip field type

	// calc key length
	it.keyLen = cstringLength(it.raw[it.offset+1:])
	fieldOffset += it.keyLen
	it.value = it.raw[it.offset+fieldOffset:]

	// calc value length
	switch t {
	case BsonTypeFloat64:
		fieldOffset += 8
	case BsonTypeString:
		fieldOffset += int(bytesToInt32(it.value)) + 4
	case BsonTypeBson:
		fallthrough
	case BsonTypeArray:
		fieldOffset += int(bytesToInt32(it.value))
	case BsonTypeBinary:
		fieldOffset += int(bytesToInt32(it.value)) + 5
	case BsonTypeObjectId:
		fieldOffset += 12
	case BsonTypeBool:
		fieldOffset += 1
	case BsonTypeDate:
		fieldOffset += 8
	case BsonTypeNull:
		// no value
	case BsonTypeRegEx:
		patternLen := cstringLength(it.value)
		optionsLen := cstringLength(it.value[patternLen:])
		fieldOffset += patternLen + optionsLen
	case BsonTypeInt32:
		fieldOffset += 4
	case BsonTypeTimestamp:
		fieldOffset += 8
	case BsonTypeInt64:
		fieldOffset += 8
	case BsonTypeMaxKey:
		// no value
	case BsonTypeMinKey:
		// no value
	default:
		panic(fmt.Sprintf("invalid bson type: %v", t))
	}
	it.elementLen = fieldOffset

	return true
}

func (it *BsonIterator) BsonType() BsonType {
	return BsonType(it.raw[it.offset])
}

func (it *BsonIterator) Name() string {
	return string(it.raw[it.offset+1 : it.offset+it.keyLen])
}

func (it *BsonIterator) Value() interface{} {
	switch it.BsonType() {
	case BsonTypeFloat64:
		return it.Float64()
	case BsonTypeString:
		return it.UTF8String()
	case BsonTypeBson:
		return it.Bson()
	case BsonTypeArray:
		return it.BsonArray()
	case BsonTypeBinary:
		return it.Binary()
	case BsonTypeObjectId:
		return it.ObjectId()
	case BsonTypeBool:
		return it.Bool()
	case BsonTypeDate:
		return it.Date()
	case BsonTypeNull:
		return nil
	case BsonTypeRegEx:
		return it.RegEx()
	case BsonTypeInt32:
		return it.Int32()
	case BsonTypeTimestamp:
		return it.Timestamp()
	case BsonTypeInt64:
		return it.Int64()
	case BsonTypeMaxKey:
		return MaxKey
	case BsonTypeMinKey:
		return MinKey
	default:
		panic(fmt.Errorf("invalid bson type: %v", it.BsonType()))
	}
}

func (it *BsonIterator) Float64() float64 {
	return bytesToFloat64(it.value)
}

func (it *BsonIterator) UTF8String() string {
	len := bytesToInt32(it.value)
	return string(it.value[4 : len+3])
}

func (it *BsonIterator) Bson() *Bson {
	len := bytesToInt32(it.value)
	return &Bson{raw: it.value[:len], finished: true}
}

func (it *BsonIterator) BsonArray() *BsonArray {
	len := bytesToInt32(it.value)
	return &BsonArray{bson: Bson{raw: it.value[:len], finished: true}}
}

func (it *BsonIterator) Binary() Binary {
	len := bytesToInt32(it.value)
	return Binary{Subtype: BinaryType(it.value[4]), Data: it.value[5 : len+5]}
}

func (it *BsonIterator) ObjectId() ObjectId {
	return ObjectId(it.value[:12])
}

func (it *BsonIterator) Bool() bool {
	return it.value[0] == 0x01
}

func (it *BsonIterator) Date() Date {
	return Date(it.Int64())
}

func (it *BsonIterator) RegEx() RegEx {
	patternLen := cstringLength(it.value)
	pattern := string(it.value[:patternLen-1])

	optionsLen := cstringLength(it.value[patternLen:])
	options := string(it.value[patternLen : patternLen+optionsLen-1])

	return RegEx{Pattern: pattern, Options: options}
}

func (it *BsonIterator) Int32() int32 {
	return bytesToInt32(it.value)
}

func (it *BsonIterator) Timestamp() Timestamp {
	inc := bytesToInt32(it.value)
	sec := bytesToInt32(it.value[4:])
	return Timestamp{Increment: inc, Second: sec}
}

func (it *BsonIterator) Int64() int64 {
	return bytesToInt64(it.value)
}
