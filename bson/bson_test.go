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

	"github.com/DavidLi2010/gobson_exp/bson"
)

func TestSingleBsonAppend(t *testing.T) {
	var tests = []struct {
		bsonType bson.BsonType
		name     string
		value    interface{}
		want     []byte
	}{
		{bson.BsonTypeFloat64, "_", 5.05, []byte("\x10\x00\x00\x00\x01_\x00333333\x14@\x00")},
		{bson.BsonTypeString, "_", "hello!", []byte("\x13\x00\x00\x00\x02_\x00\x07\x00\x00\x00hello!\x00\x00")},
		{bson.BsonTypeBinary, "_", bson.Binary{Data: []byte("abc")}, []byte("\x10\x00\x00\x00\x05_\x00\x03\x00\x00\x00\x00abc\x00")},
		{bson.BsonTypeObjectId, "_", bson.ObjectId("123456789012"), []byte("\x14\x00\x00\x00\x07_\x00123456789012\x00")},
		{bson.BsonTypeBool, "_", true, []byte("\x09\x00\x00\x00\x08_\x00\x01\x00")},
		{bson.BsonTypeDate, "_", bson.Date(12345), []byte("\x10\x00\x00\x00\x09_\x00\x39\x30\x00\x00\x00\x00\x00\x00\x00")},
		{bson.BsonTypeNull, "_", "", []byte("\x08\x00\x00\x00\x0A_\x00\x00")},
		{bson.BsonTypeRegEx, "_", bson.RegEx{Pattern: "pat", Options: "opt"}, []byte("\x10\x00\x00\x00\x0B_\x00pat\x00opt\x00\x00")},
		{bson.BsonTypeInt32, "_", int32(12345), []byte("\x0C\x00\x00\x00\x10_\x00\x39\x30\x00\x00\x00")},
		{bson.BsonTypeInt64, "_", int64(12345), []byte("\x10\x00\x00\x00\x12_\x00\x39\x30\x00\x00\x00\x00\x00\x00\x00")},
		{bson.BsonTypeTimestamp, "_", bson.Timestamp{Second: 10, Increment: 20}, []byte("\x10\x00\x00\x00\x11_\x00\x14\x00\x00\x00\x0A\x00\x00\x00\x00")},
		{bson.BsonTypeMaxKey, "_", "", []byte("\x08\x00\x00\x00\x7F_\x00\x00")},
		{bson.BsonTypeMinKey, "_", "", []byte("\x08\x00\x00\x00\xFF_\x00\x00")},
	}

	for _, test := range tests {
		doc := bson.NewBson()
		switch test.bsonType {
		case bson.BsonTypeFloat64:
			doc.AppendFloat64(test.name, test.value.(float64))
		case bson.BsonTypeString:
			doc.AppendString(test.name, test.value.(string))
		case bson.BsonTypeBinary:
			doc.AppendBinary(test.name, test.value.(bson.Binary))
		case bson.BsonTypeObjectId:
			doc.AppendObjectId(test.name, test.value.(bson.ObjectId))
		case bson.BsonTypeBool:
			doc.AppendBool(test.name, test.value.(bool))
		case bson.BsonTypeDate:
			doc.AppendDate(test.name, test.value.(bson.Date))
		case bson.BsonTypeNull:
			doc.AppendNull(test.name)
		case bson.BsonTypeRegEx:
			doc.AppendRegex(test.name, test.value.(bson.RegEx))
		case bson.BsonTypeInt32:
			doc.AppendInt32(test.name, test.value.(int32))
		case bson.BsonTypeInt64:
			doc.AppendInt64(test.name, test.value.(int64))
		case bson.BsonTypeTimestamp:
			doc.AppendTimestamp(test.name, test.value.(bson.Timestamp))
		case bson.BsonTypeMaxKey:
			doc.AppendMaxKey(test.name)
		case bson.BsonTypeMinKey:
			doc.AppendMinKey(test.name)
		default:
			t.Fatalf("invalid bson type")
		}
		doc.Finish()

		data := doc.Raw()

		if len(data) != len(test.want) {
			t.Errorf("type: %v\nexpected: %v\n  actual: %v", test.bsonType, test.want, data)
		} else {
			for i, b := range data {
				if b != test.want[i] {
					t.Errorf("type: %v\nexpected: %v\n  actual: %v", test.bsonType, test.want, data)
				}
			}
		}
	}
}

func TestBsonAppendBson(t *testing.T) {
	expected := `{"outer":"hello", "obj":{"inner":"world"}, "array":["hello world", 123.456]}`

	outer := bson.NewBson()
	outer.AppendString("outer", "hello")

	// append bson
	obj := outer.AppendBsonStart("obj")
	obj.AppendString("inner", "world")
	obj.Finish()
	outer.AppendBsonEnd()

	// append array
	array := outer.AppendArrayStart("array")
	array.AppendString("hello world")
	array.AppendFloat64(123.456)
	array.Finish()
	outer.AppendArrayEnd()

	outer.Finish()

	if expected != outer.String() {
		t.Errorf("append bson/array error, expected:%s, actual:%s", expected, outer.String())
	}
}

func TestNewBsonWithRaw(t *testing.T) {
	raw := []byte("bad bson")
	b := bson.NewBsonWithRaw(raw, bson.GetByteOrder())
	err := b.Validate()
	if err == nil {
		t.Errorf("invalid Bson.Validate()")
	}
}
