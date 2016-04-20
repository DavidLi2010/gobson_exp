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

func TestBsonIterator(t *testing.T) {
	obj := bson.NewBson()
	obj.AppendString("name", "hello")
	obj.Finish()

	array := bson.NewBsonArray()
	array.AppendString("hello")
	array.Finish()

	var tests = [...]struct {
		bsonType bson.BsonType
		name     string
		value    interface{}
	}{
		{bson.BsonTypeFloat64, "double", float64(123.45)},
		{bson.BsonTypeString, "string", "hello, bson"},
		{bson.BsonTypeBson, "bson", obj},
		{bson.BsonTypeArray, "array", array},
		{bson.BsonTypeBinary, "binary", bson.Binary{Subtype: bson.BinaryTypeGeneral, Data: []byte("hello")}},
		{bson.BsonTypeObjectId, "objectid", bson.NewObjectId()},
		{bson.BsonTypeBool, "bool", true},
		{bson.BsonTypeDate, "date", bson.Date(12345678)},
		{bson.BsonTypeNull, "null", nil},
		{bson.BsonTypeRegEx, "regex", bson.RegEx{Pattern: "/s/", Options: "g"}},
		{bson.BsonTypeInt32, "int32", int32(100)},
		{bson.BsonTypeTimestamp, "timestamp", bson.Timestamp{Second: 123, Increment: 10}},
		{bson.BsonTypeInt64, "int64", int64(5000000000)},
		{bson.BsonTypeMaxKey, "maxkey", nil},
		{bson.BsonTypeMinKey, "minkey", nil},
	}

	doc := bson.NewBson()
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		switch test.bsonType {
		case bson.BsonTypeFloat64:
			doc.AppendFloat64(test.name, test.value.(float64))
		case bson.BsonTypeString:
			doc.AppendString(test.name, test.value.(string))
		case bson.BsonTypeBson:
			doc.AppendBson(test.name, test.value.(*bson.Bson))
		case bson.BsonTypeArray:
			doc.AppendArray(test.name, test.value.(*bson.BsonArray))
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
	}
	doc.Finish()

	it := doc.Iterator()
	id := 0
	for it.Next() {
		ts := tests[id]
		tp := it.BsonType()
		name := it.Name()
		id++

		if tp != ts.bsonType {
			t.Errorf("expected type: %v, actual type: %v", ts.bsonType, tp)
		}

		if name != ts.name {
			t.Errorf("expected name: %v, len=%v; actual name: %v, len=%v",
				ts.name, len(ts.name), name, len(name))
		}

		switch tp {
		case bson.BsonTypeFloat64:
			val := it.Float64()
			if val != ts.value.(float64) {
				t.Errorf("double value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeString:
			val := it.UTF8String()
			if val != ts.value.(string) {
				t.Errorf("string value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeBson:
			val := it.Bson()
			obj := ts.value.(*bson.Bson)
			valIt := val.Iterator()
			objIt := obj.Iterator()
			for valIt.Next() && objIt.Next() {
				if valIt.BsonType() != objIt.BsonType() {
					t.Errorf("value type in doc, expected %v, actual %v", objIt.BsonType(), valIt.BsonType())
				}

				if valIt.Name() != objIt.Name() {
					t.Errorf("name in doc, expected %v, actual %v", objIt.Name(), valIt.Name())
				}
			}

			if valIt.More() || objIt.More() {
				t.Errorf("invalid it.More()")
			}
		case bson.BsonTypeArray:
			val := it.BsonArray()
			obj := ts.value.(*bson.BsonArray)
			valIt := val.Iterator()
			objIt := obj.Iterator()
			for valIt.Next() && objIt.Next() {
				if valIt.BsonType() != objIt.BsonType() {
					t.Errorf("value type in doc, expected %v, actual %v", objIt.BsonType(), valIt.BsonType())
				}

				if valIt.Name() != objIt.Name() {
					t.Errorf("name in doc, expected %v, actual %v", objIt.Name(), valIt.Name())
				}
			}
		case bson.BsonTypeBinary:
			val := it.Binary()
			expected := ts.value.(bson.Binary)
			if val.Subtype != expected.Subtype {
				t.Errorf("binary subtype, expected %v, actual %v", expected.Subtype, val.Subtype)
			}
			if string(val.Data) != string(expected.Data) {
				t.Errorf("binary data, expected %v, actual %v", expected.Data, val.Data)
			}
		case bson.BsonTypeObjectId:
			val := it.ObjectId()
			if val != ts.value.(bson.ObjectId) {
				t.Errorf("objectid value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeBool:
			val := it.Bool()
			if val != ts.value.(bool) {
				t.Errorf("bool value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeDate:
			val := it.Date()
			if val != ts.value.(bson.Date) {
				t.Errorf("date value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeNull:
			// no value
		case bson.BsonTypeRegEx:
			val := it.RegEx()
			expected := ts.value.(bson.RegEx)
			if val.Pattern != expected.Pattern || val.Options != expected.Options {
				t.Errorf("regex value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeInt32:
			val := it.Int32()
			if val != ts.value.(int32) {
				t.Errorf("int32 value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeInt64:
			val := it.Int64()
			if val != ts.value.(int64) {
				t.Errorf("int64 value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeTimestamp:
			val := it.Timestamp()
			expected := ts.value.(bson.Timestamp)
			if val.Increment != expected.Increment || val.Second != expected.Second {
				t.Errorf("timestamp value, expected %v, actual %v", ts.value, val)
			}
		case bson.BsonTypeMaxKey:
			// no value
		case bson.BsonTypeMinKey:
			// no value
		default:
			t.Fatalf("invalid bson type")
		}
	}

	if it.More() {
		t.Errorf("invalid it.More()")
	}

	if it.BsonType() != bson.BsonTypeEOD {
		t.Errorf("invalid end of iterator")
	}
}
