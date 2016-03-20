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

import "strconv"

type BsonArray struct {
	bson  Bson
	index int
}

func NewBsonArry() *BsonArray {
	return NewBsonArrayWithByteOrder(GetByteOrder())
}

func NewBsonArrayWithByteOrder(order ByteOrder) *BsonArray {
	bsonArray := &BsonArray{bson: Bson{raw: make([]byte, 0, initialBufferSize), order: order}}
	bsonArray.bson.reserveInt32()
	return bsonArray
}

func (array *BsonArray) Finish() {
	array.bson.Finish()
}

func (array *BsonArray) Raw() []byte {
	return array.bson.Raw()
}

func (array *BsonArray) AppendFloat64(value float64) {
	array.bson.AppendFloat64(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendString(value string) {
	array.bson.AppendString(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendBson(value *Bson) {
	array.bson.AppendBson(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendArray(value *BsonArray) {
	array.bson.AppendArray(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendBinary(value Binary) {
	array.bson.AppendBinary(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendObjectId(value ObjectId) {
	array.bson.AppendObjectId(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendBool(value bool) {
	array.bson.AppendBool(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendDate(value Date) {
	array.bson.AppendDate(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendNull() {
	array.bson.AppendNull(strconv.Itoa(array.index))
	array.index++
}

func (array *BsonArray) AppendRegex(pattern string, options string) {
	array.bson.AppendRegex(strconv.Itoa(array.index), pattern, options)
	array.index++
}

func (array *BsonArray) AppendInt32(value int32) {
	array.bson.AppendInt32(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendTimestamp(value Timestamp) {
	array.bson.AppendTimestamp(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendInt64(value int64) {
	array.bson.AppendInt64(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) AppendMinKey() {
	array.bson.AppendMinKey(strconv.Itoa(array.index))
	array.index++
}

func (array *BsonArray) AppendMaxKey() {
	array.bson.AppendMaxKey(strconv.Itoa(array.index))
	array.index++
}

func (array *BsonArray) Iterator() *BsonIterator {
	return NewBsonIterator(&(array.bson))
}
