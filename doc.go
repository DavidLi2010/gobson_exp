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

	// end of cstring
	eos = byte(0x00)
)

type Doc struct {
	raw []byte
}

type Date struct {
	ms  int32
	inc int32
}

type Timestamp int64

type Binary struct {
	subtype byte
	data    []byte
}

type ObjectId struct {
}

func (doc *Doc) appendDouble(name string, value float64) {

}

func (doc *Doc) appendString(name string, value string) {

}

func (doc *Doc) appendDoc(name string, value *Doc) {

}

func (doc *Doc) appendArray(name string, value []Doc) {

}

func (doc *Doc) appendBinary(name string, value Binary) {

}

func (doc *Doc) appendObjectId(name string, value ObjectId) {

}

func (doc *Doc) appendBool(name string, value bool) {

}

func (doc *Doc) appendDate(name string, value Date) {

}

func (doc *Doc) appendNull(name string) {

}

func (doc *Doc) appendRegex(name string, pattern string, options string) {

}

func (doc *Doc) appendInt32(name string, value int32) {

}

func (doc *Doc) appendTimestamp(name string, value Timestamp) {

}

func (doc *Doc) appendInt64(name string, value int64) {

}

func (doc *Doc) appendMinKey(name string) {

}

func (doc *Doc) appendMaxKey(name string) {

}
