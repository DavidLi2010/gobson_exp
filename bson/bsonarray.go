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

import (
	"bytes"
	"fmt"
	"strconv"
)

type BsonArray struct {
	bson  Bson
	index int
}

func NewBsonArray() *BsonArray {
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

func (array *BsonArray) AppendRegex(value RegEx) {
	array.bson.AppendRegex(strconv.Itoa(array.index), value)
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

func (array *BsonArray) Append(value interface{}) {
	array.bson.Append(strconv.Itoa(array.index), value)
	array.index++
}

func (array *BsonArray) Iterator() *BsonIterator {
	return NewBsonIterator(&(array.bson))
}

func (array *BsonArray) String() string {
	var err error
	buf := bytes.NewBufferString("[")
	it := array.Iterator()
	for it.Next() {
		switch it.BsonType() {
		case BsonTypeFloat64:
			_, err = fmt.Fprintf(buf, "%v", it.Float64())
		case BsonTypeString:
			_, err = fmt.Fprintf(buf, `"%s"`, it.UTF8String())
		case BsonTypeBson:
			_, err = buf.WriteString(it.Bson().String())
		case BsonTypeArray:
			_, err = buf.WriteString(it.BsonArray().String())
		case BsonTypeBinary:
			_, err = buf.WriteString(it.Binary().String())
		case BsonTypeObjectId:
			_, err = buf.WriteString(it.ObjectId().String())
		case BsonTypeBool:
			_, err = fmt.Fprintf(buf, "%v", it.Bool())
		case BsonTypeDate:
			_, err = buf.WriteString(it.Date().String())
		case BsonTypeNull:
			_, err = buf.WriteString("null")
		case BsonTypeRegEx:
			_, err = buf.WriteString(it.RegEx().String())
		case BsonTypeInt32:
			_, err = fmt.Fprintf(buf, "%d", it.Int32())
		case BsonTypeTimestamp:
			_, err = buf.WriteString(it.Timestamp().String())
		case BsonTypeInt64:
			_, err = fmt.Fprintf(buf, "%d", it.Int64())
		case BsonTypeMaxKey:
			_, err = buf.WriteString(MaxKey.String())
		case BsonTypeMinKey:
			_, err = buf.WriteString(MinKey.String())
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
			panic(fmt.Sprintf("failed to convert bson array to string: %v", err))
		}

		if it.More() {
			_, err = buf.WriteString(", ")
			if err != nil {
				panic(fmt.Sprintf("failed to convert bson array to string: %v", err))
			}
		}
	}

	_, err = buf.WriteString("]")
	if err != nil {
		panic(fmt.Sprintf("failed to convert bson array to string: %v", err))
	}

	return buf.String()
}

func (array *BsonArray) Slice() []interface{} {
	if !array.bson.finished {
		panic("the bson array is unfinished")
	}

	s := []interface{}{}

	it := array.Iterator()
	for it.Next() {
		switch it.BsonType() {
		case BsonTypeBson:
			s = append(s, it.Bson().Map())
		case BsonTypeArray:
			s = append(s, it.BsonArray().Slice())
		default:
			s = append(s, it.Value())
		}
	}

	return s
}
